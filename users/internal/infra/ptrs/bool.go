package ptrs

func Bool(bo bool) *bool {
	return &bo
}

func BoolValue(bo *bool) bool {
	if bo == nil {
		return false
	}
	return *bo
}
