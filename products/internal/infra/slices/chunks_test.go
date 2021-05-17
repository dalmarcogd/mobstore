package slices

import (
	"reflect"
	"testing"
)

func TestChunks(t *testing.T) {
	type args struct {
		slice []ChunkValue
		lim   int
	}
	tests := []struct {
		name string
		args args
		want [][]ChunkValue
	}{
		{
			name: "chunks-test1",
			args: args{
				slice: []ChunkValue{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				lim:   2,
			},
			want: [][]ChunkValue{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}},
		},
		{
			name: "chunks-test2",
			args: args{
				slice: []ChunkValue{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				lim:   1,
			},
			want: [][]ChunkValue{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}, {10}},
		},
		{
			name: "chunks-test3",
			args: args{
				slice: []ChunkValue{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				lim:   -1,
			},
			want: [][]ChunkValue{},
		},
		{
			name: "chunks-test4",
			args: args{
				slice: []ChunkValue{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				lim:   11,
			},
			want: [][]ChunkValue{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Chunks(tt.args.slice, tt.args.lim); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chunks() = %v, want %v", got, tt.want)
			}
		})
	}
}
