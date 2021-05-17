package ptrs

import "time"

func Time(it time.Time) *time.Time {
	return &it
}

func TimeValue(it *time.Time) time.Time {
	if it == nil {
		return time.Time{}
	}
	return *it
}
