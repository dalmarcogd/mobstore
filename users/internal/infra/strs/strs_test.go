package strs

import "testing"

func TestLeftPad(t *testing.T) {
	type args struct {
		s string
		n int
		r string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test-leftpad-1", args: args{
			s: "123",
			n: 10,
			r: "0",
		}, want: "0000000123"},
		{name: "test-leftpad-2", args: args{
			s: "",
			n: 10,
			r: "0",
		}, want: "0000000000"},
		{name: "test-leftpad-3", args: args{
			s: "123",
			n: 2,
			r: "0",
		}, want: "123"},
		{name: "test-leftpad-4", args: args{
			s: "123",
			n: -1,
			r: "0",
		}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LeftPad(tt.args.s, tt.args.n, tt.args.r); got != tt.want {
				t.Errorf("LeftPad() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRightPad(t *testing.T) {
	type args struct {
		s string
		n int
		r string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test-rightpad-1", args: args{
			s: "123",
			n: 10,
			r: "0",
		}, want: "1230000000"},
		{name: "test-rightpad-2", args: args{
			s: "",
			n: 10,
			r: "0",
		}, want: "0000000000"},
		{name: "test-righttpad-3", args: args{
			s: "123",
			n: 2,
			r: "0",
		}, want: "123"},
		{name: "test-righttpad-4", args: args{
			s: "123",
			n: -1,
			r: "0",
		}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RightPad(tt.args.s, tt.args.n, tt.args.r); got != tt.want {
				t.Errorf("RightPad() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpaces(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "spaces-1", args: args{
			n: 2,
		}, want: "  "},
		{name: "spaces-2", args: args{
			n: 0,
		}, want: ""},
		{name: "spaces-3", args: args{
			n: -1,
		}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Spaces(tt.args.n); got != tt.want {
				t.Errorf("Spaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToStr(t *testing.T) {
	type args struct {
		n interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "tostr-1",
			args: args{
				n: 2,
			},
			want: "2",
		},
		{
			name: "tostr-2",
			args: args{
				n: 0,
			},
			want: "0",
		},
		{
			name: "tostr-3",
			args: args{
				n: -1,
			},
			want: "-1",
		},
		{
			name: "tostr-4",
			args: args{
				n: 1.1,
			},
			want: "1.1",
		},
		{
			name: "tostr-4",
			args: args{
				n: nil,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToStr(tt.args.n); got != tt.want {
				t.Errorf("ToStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubstring(t *testing.T) {
	type args struct {
		src   string
		init  int
		final int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "substring-1", args: args{src: "teste", init: 0, final: 2}, want: "te"},
		{name: "substring-2", args: args{src: "teste", init: 1, final: 2}, want: "e"},
		{name: "substring-3", args: args{src: "teste", init: 0, final: 10}, want: "teste"},
		{name: "substring-4", args: args{src: "teste", init: 0, final: 5}, want: "teste"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Substring(tt.args.src, tt.args.init, tt.args.final); got != tt.want {
				t.Errorf("Substring() = %v, want %v", got, tt.want)
			}
		})
	}
}
