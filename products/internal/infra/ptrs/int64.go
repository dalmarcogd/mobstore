package ptrs

func Int64(it int64) *int64 {
	return &it
}

func Int64Value(it *int64) int64 {
	if it == nil {
		return 0
	}
	return *it
}
