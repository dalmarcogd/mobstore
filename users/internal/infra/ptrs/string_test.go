package ptrs

import (
	"testing"
)

func TestString(t *testing.T) {
	type args struct {
		st string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test-string-1", args: args{st: "123"}, want: "123"},
		{name: "test-string-2", args: args{st: "4444"}, want: "4444"},
		{name: "test-string-3", args: args{st: ""}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := String(tt.args.st); *got != tt.want {
				t.Errorf("String() = %v, want %v", *got, tt.want)
			}
		})
	}
}

func TestStringValue(t *testing.T) {
	type args struct {
		st *string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test-stringValue-1", args: args{st: String("123")}, want: "123"},
		{name: "test-stringValue-2", args: args{st: String("4444")}, want: "4444"},
		{name: "test-stringValue-2", args: args{st: String("")}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringValue(tt.args.st); got != tt.want {
				t.Errorf("StringValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
