package ptrs

func String(st string) *string {
	return &st
}

func StringValue(st *string) string {
	if st == nil {
		return ""
	}
	return *st
}
