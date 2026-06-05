//go:build !solution

package hogwarts

import "sort"

func GetCourseList(prereqs map[string][]string) []string {
	all := make(map[string]struct{}, len(prereqs))
	for c := range prereqs {
		all[c] = struct{}{}
	}
	for _, ps := range prereqs {
		for _, p := range ps {
			all[p] = struct{}{}
		}
	}

	indeg := make(map[string]int, len(all))
	g := make(map[string][]string, len(all))
	for c := range all {
		indeg[c] = 0
	}
	for c, ps := range prereqs {
		for _, p := range ps {
			g[p] = append(g[p], c)
			indeg[c]++
		}
	}

	zero := make([]string, 0, len(all))
	for c, d := range indeg {
		if d == 0 {
			zero = append(zero, c)
		}
	}
	sort.Strings(zero)

	order := make([]string, 0, len(all))
	for len(zero) > 0 {
		v := zero[0]
		zero = zero[1:]
		order = append(order, v)
		for _, to := range g[v] {
			indeg[to]--
			if indeg[to] == 0 {
				zero = append(zero, to)
			}
		}
		sort.Strings(zero)
	}

	if len(order) != len(all) {
		panic("cycle")
	}
	return order
}
