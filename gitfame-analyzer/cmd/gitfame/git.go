package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type BlameBlock struct {
	SHA  string
	Num  int
	Name string
	File string
}

func runGit(repo string, args ...string) ([]byte, []byte, error) {
	cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	err := cmd.Run()
	return out.Bytes(), errb.Bytes(), err
}

func RevParse(repo, rev string) (string, error) {
	stdout, stderr, err := runGit(repo, "rev-parse", "--verify", rev)
	if err != nil {
		return "", fmt.Errorf("git rev-parse %q: %w: %s", rev, err, bytes.TrimSpace(stderr))
	}
	return strings.TrimSpace(string(stdout)), nil
}

func LsTree(repo, sha string) ([]string, error) {
	stdout, stderr, err := runGit(repo, "ls-tree", "-rz", "-r", "--name-only", sha)
	if err != nil {
		return nil, fmt.Errorf("git ls-tree %q: %w: %s", sha, err, bytes.TrimSpace(stderr))
	}
	parts := bytes.Split(stdout, []byte{0})
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if len(p) != 0 {
			out = append(out, string(p))
		}
	}
	return out, nil
}

func lastTouch(repo, sha, file string, useCommitter bool) (*BlameBlock, error) {
	format := "%H%x00%an"
	if useCommitter {
		format = "%H%x00%cn"
	}
	stdout, _, err := runGit(repo, "log", "-n", "1", "--format="+format, sha, "--", file)
	if err != nil {
		return nil, nil
	}
	line := bytes.TrimSpace(stdout)
	if len(line) == 0 {
		return nil, nil
	}
	parts := bytes.SplitN(line, []byte{0}, 2)
	if len(parts) != 2 {
		return nil, nil
	}
	return &BlameBlock{
		SHA:  string(parts[0]),
		Num:  0,
		Name: string(parts[1]),
		File: file,
	}, nil
}

func BlamePorcelain(repo, sha, file string, useCommitter bool) ([]BlameBlock, error) {
	stdout, stderr, err := runGit(repo, "blame", "--line-porcelain", sha, "--", file)
	if err != nil {
		return nil, fmt.Errorf("git blame %q: %w: %s", file, err, bytes.TrimSpace(stderr))
	}
	blocks, err := parsePorcelain(stdout, useCommitter)
	if err != nil {
		return nil, err
	}
	if len(blocks) == 0 {
		if b, _ := lastTouch(repo, sha, file, useCommitter); b != nil {
			return []BlameBlock{*b}, nil
		}
	}
	for i := range blocks {
		blocks[i].File = file
	}
	return blocks, nil
}

var headerRe = regexp.MustCompile(`^[0-9a-f]{8,40}\s+\d+\s+\d+(?:\s+\d+)?$`)

func parsePorcelain(data []byte, useCommitter bool) ([]BlameBlock, error) {
	sc := bufio.NewScanner(bytes.NewReader(data))
	const max = 1024 * 1024
	buf := make([]byte, 64*1024)
	sc.Buffer(buf, max)

	var out []BlameBlock

	var curSHA string
	var curName string
	var curFile string
	var haveHdr bool

	for sc.Scan() {
		line := sc.Text()
		if headerRe.MatchString(line) {
			haveHdr = true
			i := strings.IndexByte(line, ' ')
			if i > 0 {
				curSHA = line[:i]
			}
			curName = ""
			curFile = ""
			continue
		}
		if !haveHdr {
			continue
		}
		if useCommitter {
			if strings.HasPrefix(line, "committer ") {
				curName = line[len("committer "):]
				continue
			}
		} else {
			if strings.HasPrefix(line, "author ") {
				curName = line[len("author "):]
				continue
			}
		}
		if strings.HasPrefix(line, "filename ") {
			curFile = line[len("filename "):]
			continue
		}
		if len(line) > 0 && line[0] == '\t' {
			if curSHA != "" && curName != "" && curFile != "" {
				out = append(out, BlameBlock{
					SHA:  curSHA,
					Num:  1,
					Name: curName,
					File: curFile,
				})
			}
			continue
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
