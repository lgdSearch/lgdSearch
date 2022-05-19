package trie

import "container/list"

// hotSearch use a map and a queue to get top 10 hotSearch message from today's search message
// 使用定时任务，10分钟更新一次热榜，数据通过map存储和维护，queue中存词条进入的时间和内容，那么过时间了就pop,并维护map
// 如何获取 top 10, 通过全部取出 map 中数据, 排序后取 top 10, 如果不这样，考虑我可能需要维护一个数组, 那么修改的时候再通过二分修改数组中的数据,执行冒泡操作,再取 top 10,不如第一种

type queueNode struct {
	value int
	next  *queueNode
}

type queue struct {
	size int
	head *queueNode
}

type HotSearch struct {
	searchMessage map[string]int
	list.List
}
