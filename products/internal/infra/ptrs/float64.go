package ptrs

func Float64(it float64) *float64 {
	return &it
}

func Float64Value(it *float64) float64 {
	if it == nil {
		return 0
	}
	return *it
}
