package ptrs

import (
	"testing"
)

func TestBool(t *testing.T) {
	type args struct {
		it bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "test-bool-1", args: args{it: false}, want: false},
		{name: "test-bool-2", args: args{it: true}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Bool(tt.args.it); *got != tt.want {
				t.Errorf("Bool() = %v, want %v", *got, tt.want)
			}
		})
	}
}

func TestBoolValue(t *testing.T) {
	type args struct {
		it *bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "test-boolValue-1", args: args{it: Bool(true)}, want: true},
		{name: "test-boolValue-2", args: args{it: Bool(false)}, want: false},
		{name: "test-boolValue-3", args: args{it: nil}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoolValue(tt.args.it); got != tt.want {
				t.Errorf("BoolValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
