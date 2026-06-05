package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type langMap map[string][]string

type langEntry struct {
	Name       string   `json:"name"`
	Extensions []string `json:"extensions"`
}

func normalizeExts(exts []string) []string {
	set := make(map[string]struct{}, len(exts))
	for _, raw := range exts {
		s := strings.TrimSpace(raw)
		if s == "" {
			continue
		}
		s = strings.Trim(s, `"'`)
		s = strings.ToLower(s)
		if i := strings.LastIndex(s, "."); i >= 0 {
			s = s[i:]
		} else {
			s = "." + s
		}
		if len(s) <= 1 {
			continue
		}

		set[s] = struct{}{}
	}

	out := make([]string, 0, len(set))
	for e := range set {
		out = append(out, e)
	}
	sort.Strings(out)
	return out
}

func loadLangMap() (langMap, error) {
	paths := []string{
		"configs/language_extensions.json",
		"../../configs/language_extensions.json",
		"../configs/language_extensions.json",
	}

	var data []byte
	for _, p := range paths {
		if b, e := os.ReadFile(p); e == nil {
			data = b
			break
		}
	}

	m := make(langMap)
	if len(data) > 0 {
		var arr []langEntry
		if err := json.Unmarshal(data, &arr); err == nil && len(arr) > 0 {
			for _, it := range arr {
				name := strings.ToLower(strings.TrimSpace(it.Name))
				if name == "" {
					continue
				}
				m[name] = normalizeExts(it.Extensions)
			}
		} else {
			var mp map[string][]string
			if err2 := json.Unmarshal(data, &mp); err2 == nil {
				for k, v := range mp {
					name := strings.ToLower(strings.TrimSpace(k))
					if name == "" {
						continue
					}
					m[name] = normalizeExts(v)
				}
			}
		}
	}

	defaults := map[string][]string{
		"go":       {".go"},
		"yaml":     {".yaml", ".yml"},
		"markdown": {".md"},
		"gopher":   {".gopher"},
		"c++":      {".cc", ".cpp", ".cxx", ".hh", ".hpp", ".hxx"},
		"json":     {".json"},
		"text":     {".txt"},
		"proto":    {".proto"},
	}
	for k, v := range defaults {
		if _, ok := m[k]; !ok {
			m[k] = v
		}
	}

	return m, nil
}

func ResolveExts(exts, langs []string) ([]string, error) {
	m, _ := loadLangMap()
	set := make(map[string]struct{})

	for _, e := range exts {
		e = strings.ToLower(strings.TrimSpace(e))
		if e == "" {
			continue
		}
		if !strings.HasPrefix(e, ".") {
			e = "." + e
		}
		set[e] = struct{}{}
	}

	for _, l := range langs {
		ls := strings.ToLower(strings.TrimSpace(l))
		if ls == "" {
			continue
		}
		if mm, ok := m[ls]; ok {
			for _, e := range mm {
				set[e] = struct{}{}
			}
		} else {
			fmt.Fprintf(os.Stderr, "warning: unknown language: %s\n", l)
		}
	}

	out := make([]string, 0, len(set))
	for e := range set {
		out = append(out, e)
	}
	return out, nil
}

func filterByExt(paths, exts []string) []string {
	if len(exts) == 0 {
		return paths
	}
	allowed := make(map[string]struct{}, len(exts))
	for _, e := range exts {
		allowed[strings.ToLower(e)] = struct{}{}
	}
	var res []string
	for _, p := range paths {
		ext := strings.ToLower(filepath.Ext(p))
		if _, ok := allowed[ext]; ok {
			res = append(res, p)
		}
	}
	return res
}

func filterRestrict(paths, globs []string) []string {
	if len(globs) == 0 {
		return paths
	}
	var res []string
	for _, p := range paths {
		for _, g := range globs {
			if ok, _ := filepath.Match(g, p); ok {
				res = append(res, p)
				break
			}
		}
	}
	return res
}

func filterExclude(paths, globs []string) []string {
	if len(globs) == 0 {
		return paths
	}
	var res []string
next:
	for _, p := range paths {
		for _, g := range globs {
			if ok, _ := filepath.Match(g, p); ok {
				continue next
			}
		}
		res = append(res, p)
	}
	return res
}

func ApplyFilters(paths []string, cfg Config) ([]string, error) {
	exts, err := ResolveExts(cfg.Extensions, cfg.Languages)
	if err != nil {
		return nil, err
	}
	if len(cfg.Languages) > 0 && len(exts) == 0 {
		return []string{}, nil
	}
	out := filterByExt(paths, exts)
	out = filterRestrict(out, cfg.RestrictGlobs)
	out = filterExclude(out, cfg.ExcludeGlobs)
	return out, nil
}
