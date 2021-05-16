package slices

type ChunkValue interface{}

func Chunks(slice []ChunkValue, lim int) [][]ChunkValue {
	if lim < 0 {
		return make([][]ChunkValue, 0)
	}
	var chunk []ChunkValue
	chunks := make([][]ChunkValue, 0, len(slice)/lim+1)
	for len(slice) >= lim {
		chunk, slice = slice[:lim], slice[lim:]
		chunks = append(chunks, chunk)
	}
	if len(slice) > 0 {
		chunks = append(chunks, slice[:])
	}
	return chunks
}
