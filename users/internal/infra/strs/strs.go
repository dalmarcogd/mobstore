package strs

import (
	"fmt"
	"strings"
)

func LeftPad(s string, n int, r string) string {
	if n < 0 {
		return ""
	}
	if len(s) > n {
		return s
	}
	return strings.Repeat(r, n-len(s)) + s
}

func RightPad(s string, n int, r string) string {
	if n < 0 {
		return ""
	}
	if len(s) > n {
		return s
	}
	return s + strings.Repeat(r, n-len(s))
}

func Spaces(n int) string {
	if n < 0 {
		return ""
	}
	return strings.Repeat(" ", n)
}

func ToStr(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func Substring(src string, init, final int) string {
	s := len(src)
	if final > s {
		final = s
	}
	return src[init:final]
}
