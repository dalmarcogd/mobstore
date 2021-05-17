package slices

import "strconv"

func ParseInt64ToStr(ns []int64) []string {
	ss := make([]string, len(ns))
	for i, n := range ns {
		ss[i] = strconv.FormatInt(n, 10)
	}
	return ss
}

func ParseStrToInt64(ns []string) []int64 {
	ss := make([]int64, 0)
	for _, n := range ns {
		parsedInt, err := strconv.ParseInt(n, 10, 64)
		if err == nil {
			ss = append(ss, parsedInt)
		}
	}
	return ss
}
