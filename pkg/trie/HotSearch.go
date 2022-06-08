package trie

import (
	"container/heap"
	"lgdSearch/pkg/utils"
	"os"
	"runtime"
	"sync"
	"time"
)

const HotSearchFileName = "HotSearch.txt"

// hotSearch use a map and a queue to get top 10 hotSearch message from today's search message
// 使用定时任务，10分钟更新一次热榜，数据通过map存储和维护，queue中存词条进入的时间和内容，那么过时间了就pop,并维护map
type queueNode struct {
	TimeMessage time.Time
	Text        string
	Next        *queueNode
}

type queue struct {
	size int
	head *queueNode
	end  *queueNode
	sync.Mutex
}

func (q *queue) Size() int {
	return q.size
}

func (q *queue) Push(node *queueNode) {
	q.Lock()
	defer q.Unlock()
	q.size += 1
	if q.size == 1 {
		q.head = node
	} else {
		q.end.Next = node
	}
	q.end = node
}

func (q *queue) Pop() *queueNode {
	q.Lock()
	defer q.Unlock()
	if q.size == 0 {
		return q.head
	}
	top := q.head
	q.head = q.head.Next
	q.size -= 1
	if q.size == 0 { // 可以不写，其实是不影响的，不过考虑严谨性还是写了
		q.end = nil
	}
	return top
}

func (q *queue) Top() *queueNode {
	return q.head
}

var MessageChan chan string

type HotSearch struct {
	searchMessage  map[string]int // 存热点数据的搜索次数
	searchQueue    *queue         // queue 维护map中的数据
	hotSearchArray []*HotSearchMessage
	sync.Mutex
}

var MyHotSearch *HotSearch

func InitHotSearch(filepath string) {
	GetHotSearch()
	GetHotSearch().Load(filepath)
	go GetHotSearch().InQueueExec()
	go GetHotSearch().OutQueueExec()
	go GetHotSearch().AutoReGetArray(filepath)
	GetHotSearch().ReGetArray()
}

func GetHotSearch() *HotSearch {
	if MyHotSearch == nil {
		MyHotSearch = &HotSearch{
			searchMessage:  make(map[string]int),
			searchQueue:    &queue{},
			hotSearchArray: make([]*HotSearchMessage, 10),
		}

		MessageChan = make(chan string, 1000)
	}
	return MyHotSearch
}

func (hot *HotSearch) Array() []*HotSearchMessage {
	return hot.hotSearchArray
}

type HotSearchMessage struct {
	Text string `json:"text,omitempty"`
	Num  int    `json:"num,omitempty"`
}

func (node1 *HotSearchMessage) compare(node2 *HotSearchMessage) bool {
	return node1.Num < node2.Num
}

type hotHeap []*HotSearchMessage

func (h hotHeap) Len() int { return len(h) }

func (h hotHeap) Less(i, j int) bool {
	return h[i].compare(h[j])
}

func (h hotHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *hotHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*HotSearchMessage))
}

func (h *hotHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (hot *HotSearch) OutQueueExec() {
	for {
		// 这是考虑时效性的写法
		if hot.searchQueue.Size() != 0 &&
			time.Now().Sub(hot.searchQueue.Top().TimeMessage).Hours() > 24. {
			node := hot.searchQueue.Pop()
			hot.searchMessage[node.Text]--
		}
		// 如果考虑数据量的写法
		// if hot.searchQueue.Size() > 10000000 {}
	}
}

func SendText(text string) {
	MessageChan <- text
}

func (hot *HotSearch) InQueueExec() {
	for {
		text := <-MessageChan
		hot.searchQueue.Push(&queueNode{TimeMessage: time.Now(), Text: text})
		hot.searchMessage[text]++
	}
}

// AutoReGetArray 120S 自动更新一次
func (hot *HotSearch) AutoReGetArray(filepath string) {
	ticker := time.NewTicker(time.Second * 120)
	head := hot.searchQueue.head
	size := hot.searchQueue.Size()

	for {
		<-ticker.C
		if head == hot.searchQueue.head && size == hot.searchQueue.Size() {
			continue
		}
		head = hot.searchQueue.head
		size = hot.searchQueue.Size()
		hot.ReGetArray()
		hot.Flush(filepath)

		//定时GC
		runtime.GC()
	}
}

func (hot *HotSearch) ReGetArray() {
	hot.Lock() // lock 所以需要一个 chan 来暂存这些 query
	defer hot.Unlock()

	minHeap := hotHeap{} // string 的 minheap
	for k, v := range hot.searchMessage {
		if minHeap.Len() < 10 {
			heap.Push(&minHeap, &HotSearchMessage{Text: k, Num: v})
		} else if minHeap[0].Num < v {
			heap.Pop(&minHeap)
			heap.Push(&minHeap, &HotSearchMessage{Text: k, Num: v})
		}
	}

	hotMessages := make([]HotSearchMessage, 0)
	for minHeap.Len() > 0 {
		val := heap.Pop(&minHeap)
		hotMessages = append(hotMessages, *(val.(*HotSearchMessage)))
	} // 这是小顶堆，需要反转
	for i, j := 0, len(hotMessages)-1; i < j; i, j = i+1, j-1 {
		hotMessages[i], hotMessages[j] = hotMessages[j], hotMessages[i]
	}

	for index, val := range hotMessages {
		hot.hotSearchArray[index] = &HotSearchMessage{Text: val.Text, Num: val.Num}
	}
	for i := len(hotMessages); i < 10; i++ {
		hot.hotSearchArray[i] = nil
	}
}

func (hot *HotSearch) Flush(filepath string) {
	data := make([]queueNode, 0)
	for head := hot.searchQueue.head; head != nil; head = head.Next {
		data = append(data, queueNode{Text: head.Text, TimeMessage: head.TimeMessage, Next: nil})
		if head == head.Next {
			break
		}
	}
	file, _ := os.OpenFile(filepath, os.O_CREATE|os.O_TRUNC, 0600)
	file.Close()

	utils.Write(&data, filepath)
}

func (hot *HotSearch) Load(filepath string) {
	if filepath == "" {
		filepath = "./pkg/data/HotSearch.txt"
	}
	data := make([]queueNode, 0)
	utils.Read(&data, filepath)
	for _, val := range data {
		hot.searchQueue.Push(&val)
		hot.searchMessage[val.Text]++
	}
}
