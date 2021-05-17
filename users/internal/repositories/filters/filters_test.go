package filters

import (
	"testing"

	"github.com/dalmarcogd/mobstore/users/internal/infra/ptrs"
)

func TestGetFilter(t *testing.T) {
	type args struct {
		ts interface{}
	}

	type test struct {
		name string
		args args
		want map[string]interface{}
	}
	type testValue struct {
		Name   string `filter:"Name"`
		Value1 int64  `filter:"Value"`
		Value2 *int64 `filter:"Amount"`
	}

	var tests []test
	tests = append(tests, test{name: "GetFilters-1", args: args{ts: testValue{Name: "234"}}, want: map[string]interface{}{"Name": "234"}})
	tests = append(tests, test{name: "GetFilters-2", args: args{ts: testValue{Name: "234", Value1: 999}}, want: map[string]interface{}{"Name": "234", "Value": int64(999)}})
	tests = append(tests, test{name: "GetFilters-3", args: args{ts: testValue{Name: "234", Value2: ptrs.Int64(1)}}, want: map[string]interface{}{"Name": "234", "Amount": int64(1)}})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFilters(tt.args.ts); got != nil {
				for s, vw := range tt.want {
					if v, ok := got[s]; ok {
						if v == vw {
							delete(got, s)
						}
					}
				}
				if len(got) > 0 {
					t.Errorf("GetFilters() unexpected value = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
