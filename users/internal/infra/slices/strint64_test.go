package slices

import (
	"reflect"
	"testing"
)

func TestParseSliceInt64ToSliceStr(t *testing.T) {
	tests := []struct {
		name string
		args []int64
		want []string
	}{
		{name: "test1", args: []int64{10, 11, 112}, want: []string{"10", "11", "112"}},
		{name: "test2", args: []int64{}, want: []string{}},
		{name: "test3", args: []int64{10}, want: []string{"10"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseInt64ToStr(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%v ParseInt64ToStr() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestParseSliceStrToSliceInt64(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []int64
	}{
		{name: "test1", args: []string{"10", "11", "112", "10.123"}, want: []int64{10, 11, 112}},
		{name: "test2", args: []string{}, want: []int64{}},
		{name: "test3", args: []string{"10"}, want: []int64{10}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseStrToInt64(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%v ParseStrToInt64() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
