package reqgetter

import (
	"bufio"
	"bytes"
	"context"
	"net/http"
	"testing"
)

func TestGetCid(t *testing.T) {
	type args struct {
		req *http.Request
	}
	re1, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", bufio.NewReader(bytes.NewReader([]byte("t"))))
	if err != nil {
		t.Error(err)
	}
	re1.Header.Add("x-cid", "test-cid-1")
	re2, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", bufio.NewReader(bytes.NewReader([]byte("t"))))
	if err != nil {
		t.Error(err)
	}
	re2.Header.Add("x-cid", "test-cid-2")
	re3, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", bufio.NewReader(bytes.NewReader([]byte("t"))))
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name string
		args args
		want func() string
	}{
		{name: "GetCid-1", args: args{req: re1}, want: func() string { return re1.Header.Get("x-cid") }},
		{name: "GetCid-2", args: args{req: re2}, want: func() string { return re2.Header.Get("x-cid") }},
		{name: "GetCid-3", args: args{req: re3}, want: func() string { return re3.Header.Get("x-cid") }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v string
			got := GetCid(tt.args.req)
			if got != nil {
				v = *got
			}
			if v != tt.want() {
				t.Errorf("GetCid() = %v, want %v", v, tt.want())
			}
		})
	}
}
