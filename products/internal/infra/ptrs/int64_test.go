package ptrs

import (
	"testing"
)

func TestInt64(t *testing.T) {
	type args struct {
		it int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{name: "test-int64-1", args: args{it: int64(0)}, want: int64(0)},
		{name: "test-int64-2", args: args{it: int64(10)}, want: int64(10)},
		{name: "test-int64-2", args: args{it: int64(-1)}, want: int64(-1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int64(tt.args.it); *got != tt.want {
				t.Errorf("Int64() = %v, want %v", *got, tt.want)
			}
		})
	}
}

func TestInt64Value(t *testing.T) {
	type args struct {
		it *int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{name: "test-int64Value-1", args: args{it: Int64(0)}, want: 0},
		{name: "test-int64Value-2", args: args{it: Int64(10)}, want: 10},
		{name: "test-int64Value-2", args: args{it: Int64(-1)}, want: -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int64Value(tt.args.it); got != tt.want {
				t.Errorf("Int64Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
