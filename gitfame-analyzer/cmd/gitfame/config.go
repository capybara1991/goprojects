package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"runtime"
	"strings"
)

type OrderBy int

const (
	OrderByLines OrderBy = iota
	OrderByCommits
	OrderByFiles
)

type Format int

const (
	FormatTabular Format = iota
	FormatCSV
	FormatJSON
	FormatJSONL
)

type Config struct {
	Repository, Revision string
	UseCommitter         bool
	OrderBy              OrderBy
	Format               Format
	Extensions           []string
	Languages            []string
	RestrictGlobs        []string
	ExcludeGlobs         []string
	Progress             bool
	Workers              int
}

func ParseFlags() (Config, error) {
	repository := pflag.String("repository", ".", "")
	revision := pflag.String("revision", "HEAD", "")
	orderByStr := pflag.String("order-by", "lines", "")
	useCommitter := pflag.Bool("use-committer", false, "")
	formatStr := pflag.String("format", "tabular", "")

	extensions := pflag.String("extensions", "", "")
	languages := pflag.String("languages", "", "")
	restrictTo := pflag.String("restrict-to", "", "")
	exclude := pflag.String("exclude", "", "")

	progress := pflag.Bool("progress", false, "")
	workers := pflag.Int("workers", 0, "")

	pflag.Parse()

	if fi, err := os.Stat(*repository); err != nil || !fi.IsDir() {
		return Config{}, fmt.Errorf("invalid --repository %q: %v", *repository, err)
	}
	var ob OrderBy
	switch strings.ToLower(strings.TrimSpace(*orderByStr)) {
	case "lines":
		ob = OrderByLines
	case "commits":
		ob = OrderByCommits
	case "files":
		ob = OrderByFiles
	default:
		return Config{}, fmt.Errorf("invalid --order-by %q (use: lines|commits|files)", *orderByStr)
	}
	var fm Format
	switch strings.ToLower(strings.TrimSpace(*formatStr)) {
	case "tabular":
		fm = FormatTabular
	case "csv":
		fm = FormatCSV
	case "json":
		fm = FormatJSON
	case "json-lines", "jsonl", "jsonlines":
		fm = FormatJSONL
	default:
		return Config{}, fmt.Errorf("invalid --format %q (use: tabular|csv|json|json-lines)", *formatStr)
	}

	exts := splitCSV(*extensions)
	exts = normalizeExts(exts)
	langs := splitCSV(*languages)
	restrict := splitCSV(*restrictTo)
	excl := splitCSV(*exclude)

	nw := *workers
	if nw <= 0 {
		if n := runtime.NumCPU(); n < 8 {
			nw = n
		} else {
			nw = 8
		}
	}

	cfg := Config{
		Repository:    *repository,
		Revision:      *revision,
		UseCommitter:  *useCommitter,
		OrderBy:       ob,
		Format:        fm,
		Extensions:    exts,
		Languages:     langs,
		RestrictGlobs: restrict,
		ExcludeGlobs:  excl,
		Progress:      *progress,
		Workers:       nw,
	}
	return cfg, nil
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	seen := make(map[string]struct{}, len(parts))
	out := make([]string, 0, len(parts))

	for _, part := range parts {
		p := strings.TrimSpace(part)
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}

	return out
}
