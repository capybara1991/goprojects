package main

import "sort"

type Row struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}

type agg struct {
	lines   int
	commits map[string]struct{}
	files   map[string]struct{}
}

type Acc struct{ m map[string]*agg }

func NewAcc() *Acc { return &Acc{m: make(map[string]*agg)} }

func (a *Acc) Add(b BlameBlock) {
	x, ok := a.m[b.Name]
	if !ok {
		x = &agg{commits: make(map[string]struct{}), files: make(map[string]struct{})}
		a.m[b.Name] = x
	}
	if b.Num > 0 {
		x.lines += b.Num
	}
	if b.SHA != "" {
		x.commits[b.SHA] = struct{}{}
	}
	if b.File != "" {
		x.files[b.File] = struct{}{}
	}
}

func (a *Acc) Rows() []Row {
	rows := make([]Row, 0, len(a.m))
	for name, x := range a.m {
		rows = append(rows, Row{
			Name:    name,
			Lines:   x.lines,
			Commits: len(x.commits),
			Files:   len(x.files),
		})
	}
	return rows
}

func SortRows(rows []Row, by OrderBy) []Row {
	sort.Slice(rows, func(i, j int) bool {
		switch by {
		case OrderByCommits:
			if rows[i].Commits != rows[j].Commits {
				return rows[i].Commits > rows[j].Commits
			}
			if rows[i].Lines != rows[j].Lines {
				return rows[i].Lines > rows[j].Lines
			}
			if rows[i].Files != rows[j].Files {
				return rows[i].Files > rows[j].Files
			}
		case OrderByFiles:
			if rows[i].Files != rows[j].Files {
				return rows[i].Files > rows[j].Files
			}
			if rows[i].Lines != rows[j].Lines {
				return rows[i].Lines > rows[j].Lines
			}
			if rows[i].Commits != rows[j].Commits {
				return rows[i].Commits > rows[j].Commits
			}
		default:
			if rows[i].Lines != rows[j].Lines {
				return rows[i].Lines > rows[j].Lines
			}
			if rows[i].Commits != rows[j].Commits {
				return rows[i].Commits > rows[j].Commits
			}
			if rows[i].Files != rows[j].Files {
				return rows[i].Files > rows[j].Files
			}
		}
		return rows[i].Name < rows[j].Name
	})
	return rows
}
