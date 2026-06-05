package main

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"strconv"
	"text/tabwriter"
)

func PrintTabular(rows []Row) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	_, _ = w.Write([]byte("Name\tLines\tCommits\tFiles\n"))
	for _, r := range rows {
		_, _ = w.Write([]byte(
			r.Name + "\t" + itoa(r.Lines) + "\t" + itoa(r.Commits) + "\t" + itoa(r.Files) + "\n",
		))
	}
	return w.Flush()
}

func PrintCSV(rows []Row) error {
	w := csv.NewWriter(os.Stdout)
	_ = w.Write([]string{"Name", "Lines", "Commits", "Files"})
	for _, r := range rows {
		_ = w.Write([]string{r.Name, itoa(r.Lines), itoa(r.Commits), itoa(r.Files)})
	}
	w.Flush()
	return w.Error()
}

func PrintJSON(rows []Row) error {
	enc := json.NewEncoder(os.Stdout)
	return enc.Encode(rows)
}

func PrintJSONL(rows []Row) error {
	enc := json.NewEncoder(os.Stdout)
	for _, r := range rows {
		if err := enc.Encode(r); err != nil {
			return err
		}
	}
	return nil
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
