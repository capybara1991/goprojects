package main

import (
	"context"
)

func Run(ctx context.Context, cfg Config) error {
	sha, err := RevParse(cfg.Repository, cfg.Revision)
	if err != nil {
		return err
	}
	paths, err := LsTree(cfg.Repository, sha)
	if err != nil {
		return err
	}
	paths, err = ApplyFilters(paths, cfg)
	if err != nil {
		return err
	}
	acc := NewAcc()
	for _, p := range paths {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		blocks, err := BlamePorcelain(cfg.Repository, sha, p, cfg.UseCommitter)
		if err != nil {
			return err
		}
		for _, b := range blocks {
			acc.Add(b)
		}
	}
	rows := SortRows(acc.Rows(), cfg.OrderBy)
	switch cfg.Format {
	case FormatTabular:
		return PrintTabular(rows)
	case FormatCSV:
		return PrintCSV(rows)
	case FormatJSON:
		return PrintJSON(rows)
	case FormatJSONL:
		return PrintJSONL(rows)
	}
	return nil
}
