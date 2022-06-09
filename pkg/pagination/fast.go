package pagination

import (
	"lgdSearch/pkg/models"
	"log"
	"sort"
	"time"
)

type ScoreSlice []models.SliceItem

func (x ScoreSlice) Len() int {
	return len(x)
}
func (x ScoreSlice) Less(i, j int) bool {
	return x[i].Score < x[j].Score
}
func (x ScoreSlice) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

type FastSort struct {
	ScoreMap map[uint32]float32
	Data     []*models.SliceItem
}

// Add :Count the scores of all documents
func (f *FastSort) Add(values []*models.SliceItem) {
	if values == nil {
		return
	}
	for _, item := range values {
		f.ScoreMap[item.Id] += item.Score
	}
}

// Count 获取数量
func (f *FastSort) Count() int {
	return len(f.ScoreMap)
}

// GetAll 获取按 score 排序后的结果集
func (f *FastSort) GetAll() []models.SliceItem {
	delete(f.ScoreMap, 0) // 如果有 0 去掉

	var result = make([]models.SliceItem, len(f.ScoreMap))

	_time := time.Now()
	index := 0
	for key, value := range f.ScoreMap {
		if key == 0 {
			continue
		}
		result[index] = models.SliceItem{Id: key, Score: value}
		index++
	}
	log.Println("fastSort: 取数据耗时", time.Since(_time))

	_time = time.Now()
	sort.Sort(sort.Reverse(ScoreSlice(result)))
	log.Println("fastSort: 排序耗时", time.Since(_time))

	return result
}
