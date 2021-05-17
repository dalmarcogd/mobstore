package ptrs

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	type args struct {
		it time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{name: "test-int64-1", args: args{it: time.Time{}.Add(time.Hour)}, want: time.Time{}.Add(time.Hour)},
		{name: "test-int64-2", args: args{it: time.Time{}.Add(time.Duration(time.Month(1)))}, want: time.Time{}.Add(time.Duration(time.Month(1)))},
		{name: "test-int64-2", args: args{it: time.Time{}}, want: time.Time{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Time(tt.args.it); *got != tt.want {
				t.Errorf("Time() = %v, want %v", *got, tt.want)
			}
		})
	}
}

func TestTimeValue(t *testing.T) {
	type args struct {
		it *time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{name: "test-timeValue-1", args: args{it: Time(time.Time{}.Add(time.Hour))}, want: time.Time{}.Add(time.Hour)},
		{name: "test-timeValue-2", args: args{it: Time(time.Time{}.Add(time.Duration(time.Month(1))))}, want: time.Time{}.Add(time.Duration(time.Month(1)))},
		{name: "test-timeValue-2", args: args{it: Time(time.Time{})}, want: time.Time{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeValue(tt.args.it); got != tt.want {
				t.Errorf("TimeValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
