//go:build !solution

package hotelbusiness

import "sort"

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {
	if len(guests) == 0 {
		return nil
	}
	delta := make(map[int]int, len(guests)*2)
	for _, g := range guests {
		delta[g.CheckInDate]++
		delta[g.CheckOutDate]--
	}
	dates := make([]int, 0, len(delta))
	for d := range delta {
		dates = append(dates, d)
	}
	sort.Ints(dates)

	res := make([]Load, 0, len(dates))
	prev := 0
	cur := 0
	for _, d := range dates {
		cur += delta[d]
		if cur != prev {
			res = append(res, Load{StartDate: d, GuestCount: cur})
		}
		prev = cur
	}
	return res
}
