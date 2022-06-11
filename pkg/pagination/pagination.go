package pagination

import (
	"math"
)

type Pagination struct {
	Limit int //限制大小

	PageCount int //总页数
	Total     int //总数据量
}

func (p *Pagination) Init(limit int, total int) {
	p.Limit = limit

	p.Total = total

	pageCount := math.Ceil(float64(total) / float64(limit))
	p.PageCount = int(pageCount)

}

func (p *Pagination) GetPage(page int) (s int, e int) {
	//获取指定页数的数据
	if page > p.PageCount { // 范围超出，不返回
		return -1, -1
	}
	if page < 0 {
		page = 1
	}

	//从 0 开始
	page -= 1

	//计算起始位置
	start := page * p.Limit
	end := (start - 1) + p.Limit

	if start >= p.Total {
		//return 0, p.Total - 1
		return -1, -1
	}
	if end >= p.Total {
		end = p.Total - 1
	}

	return start, end
}
