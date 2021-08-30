package node

import (
	"github.com/OnFinality-io/onf-cli/pkg/utils"
	"math"
	"testing"
)

func TestChunk(t *testing.T) {

	var arr []int
	for i := 0; i < 3; i++ {
		arr = append(arr, i)
	}
	l := len(arr)
	percent := 1
	chunkSize := float64(l) * float64(utils.Min(100, percent)) / 100
	size := int(math.Floor(chunkSize))
	if size == 0 {
		size = 1
	}

	t.Log(chunkSize, size)
	newArr := chunk(arr, size)
	t.Log(newArr)
}
func chunk(slice []int, chunkSize int) [][]int {
	var chunks [][]int
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
