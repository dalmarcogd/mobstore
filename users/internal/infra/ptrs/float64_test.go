package ptrs

import (
	"testing"
)

func TestFloat64(t *testing.T) {
	type args struct {
		it float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{name: "test-Float64-1", args: args{it: float64(0)}, want: float64(0)},
		{name: "test-Float64-2", args: args{it: float64(10)}, want: float64(10)},
		{name: "test-Float64-2", args: args{it: float64(-1)}, want: float64(-1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Float64(tt.args.it); *got != tt.want {
				t.Errorf("Float64() = %v, want %v", *got, tt.want)
			}
		})
	}
}

func TestFloat64Value(t *testing.T) {
	type args struct {
		it *float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{name: "test-Float64Value-1", args: args{it: Float64(0)}, want: 0},
		{name: "test-Float64Value-2", args: args{it: Float64(10)}, want: 10},
		{name: "test-Float64Value-2", args: args{it: Float64(-1)}, want: -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Float64Value(tt.args.it); got != tt.want {
				t.Errorf("Float64Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
