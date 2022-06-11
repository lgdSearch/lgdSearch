package pagination

import (
	"testing"
)

func TestPagination_GetPage(t *testing.T) {
	var data []int64
	for i := 0; i < 100; i++ {
		data = append(data, int64(i))
	}
}
