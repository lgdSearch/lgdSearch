package pagination

import (
	"lgdSearch/pkg/logger"
	"lgdSearch/pkg/models"
	"sort"
	"time"
)

type ScoreSlice []DocItem

func (x ScoreSlice) Len() int {
	return len(x)
}

func (x ScoreSlice) Less(i, j int) bool {
	if x[i].Count != x[j].Count {
		return x[i].Count < x[j].Count
	} else {
		return x[i].Score < x[j].Score
	}
}

func (x ScoreSlice) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

type DocScore struct {
	// 分词总分数
	score float32

	// 符合分词个数
	count int
}

type DocItem struct {
	Id    uint32
	Score float32
	Count int
}

type FastSort struct {
	ScoreMap map[uint32]*DocScore
}

// Add :Count the scores of all documents
func (f *FastSort) Add(values []*models.SliceItem) {
	if values == nil {
		return
	}
	for _, item := range values {
		_, ok := f.ScoreMap[item.Id]
		if ok {
			f.ScoreMap[item.Id].score += item.Score
			f.ScoreMap[item.Id].count += 1
		} else {
			f.ScoreMap[item.Id] = &DocScore{score: item.Score, count: 1}
		}
	}
}

// Count 获取数量
func (f *FastSort) Count() int {
	return len(f.ScoreMap)
}

// GetAll 获取按 score 排序后的结果集
func (f *FastSort) GetAll() []DocItem {
	delete(f.ScoreMap, 0) // 如果有 0 去掉

	var result = make([]DocItem, len(f.ScoreMap))

	_time := time.Now()
	index := 0
	for key, value := range f.ScoreMap {
		if key == 0 {
			continue
		}
		result[index] = DocItem{Id: key, Score: value.score, Count: value.count}
		index++
	}
	logger.Logger.Infoln("fastSort: 取数据耗时", time.Since(_time))

	_time = time.Now()
	sort.Sort(sort.Reverse(ScoreSlice(result)))
	logger.Logger.Infoln("fastSort: 排序耗时", time.Since(_time))

	return result
}
