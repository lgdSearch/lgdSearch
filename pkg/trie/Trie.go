// Copyright (c) 2022, wanshe li. All rights reserved

package trie

import (
	"container/heap"
	"fmt"
	"sync"
)

// Node is a single element within the Trie
type Node struct {
	father *Node // 双向关系
	data   rune
	child  map[rune]*Node

	size  int32 // 统计子树完整串数量
	count int32 // 累计查找次数(用于排序)，这个是以这个词为结尾的句子个数 : int32 = rune
	max   int32 // 维护子树中最大 count, 部分节点可能 count = 0, 但是 max 有值，这表明它是一个非结尾，但是它下面有结尾
}

// get a new Node
func newNode() *Node {
	return &Node{child: make(map[rune]*Node)}
}

// string returns a string describe Node
func (n *Node) string() string {
	return fmt.Sprintf("Node: {data: %s, count: %d, child: %v}", string(n.data), n.count, n.child)
}

// Trie holds elements of the Trie tree.
type Trie struct {
	size int
	root *Node
	sync.Mutex
}

// NewTrie get a new Trie
func NewTrie() *Trie {
	return &Trie{root: newNode()}
}

// string returns a string representation of container
func (t *Trie) string() string {
	str := "Trie\n"
	if t.root != nil {
		output(t.root, &str)
	}
	return str
}

func output(rt *Node, str *string) {

}

// get new root Node and set size to zero
func (t *Trie) init() {
	t.size = 0
	t.root = newNode()
}

// change root Node to nil
func (t *Trie) clear() {
	t.size = 0
	t.root = nil // 自动 GC
}

// InsertBytes insert bytes into Trie
func (t *Trie) InsertBytes(bytes []byte) {
	str := string(bytes)
	t.InsertString(str)
}

// InsertString insert string into Trie
func (t *Trie) InsertString(str string) {
	runes := []rune(str)
	t.InsertRunes(runes)
}

// InsertRunes insert runes into Trie
func (t *Trie) InsertRunes(runes []rune) {
	t.insertRunesWithCount(runes, 1)
}

// insert runes into Trie with count
func (t *Trie) insertRunesWithCount(runes []rune, count int32) {
	if len(runes) == 0 {
		return
	}
	now := t.root
	stack := make([]*Node, len(runes))
	flag := false
	for _, val := range runes {
		nxt, ok := now.child[val]
		if !ok {
			nxt = newNode()
			nxt.data = val
			now.child[val] = nxt
			nxt.father = now
			t.size += 1
		}
		// 递归维护 max, size
		defer func(now *Node, nxt *Node, flag *bool) {
			if nxt.max > now.max {
				now.max = nxt.max
			}
			if *flag {
				now.size += 1
			}
		}(now, nxt, &flag)
		now = nxt
		stack = append(stack, now)
	}
	if now.count == 0 { // 第一次成为完整串, 贡献一个 size
		now.size += 1
		flag = true
	}
	now.count += count // 此句出现count次
}

// find the last character in string from Trie.
// If not find, return nil
func (t *Trie) findByString(str string) (*Node, int32) {
	runes := []rune(str)
	return t.findByRunes(runes)
}

// find the last character in bytes from Trie.
// If not find, return nil
func (t *Trie) findByBytes(bytes []byte) (*Node, int32) {
	str := string(bytes)
	return t.findByString(str)
}

// find the last character in runes from Trie
// If not find, return nil
func (t *Trie) findByRunes(runes []rune) (*Node, int32) {
	if len(runes) == 0 {
		return nil, 0
	}
	//defer t.InsertRunes(runes) // 查找之后插入此查询
	now := t.root
	deep := int32(0)
	for _, val := range runes {
		nxt, ok := now.child[val]
		if !ok {
			return now, deep
		}
		deep += 1
		now = nxt
	}
	return now, deep
}

// Query is a heap sizeof 10, save related search
type Query struct {
	heap            Heap
	addHeapNodeChan chan *heapNode
	wait            sync.WaitGroup
}

// 获取前缀
func getPrefix(node *Node) *string {
	str := ""
	for node.father != nil {
		defer func(str *string, node *Node) {
			*str += string(node.data)
		}(&str, node)
		node = node.father
	}
	return &str
}

// Search 查找相关搜索词条
func (t *Trie) Search(runes []rune) []string {
	//defer t.InsertRunes(runes) // 查找后插入数据
	query := &Query{addHeapNodeChan: make(chan *heapNode, 10), heap: Heap{}}
	node, deep := t.findByRunes(runes)

	h := MaxHeap{}
	heap.Push(&h, &heapNode{node: node, deep: deep})
	for i := deep - 1; i > 0; i-- {
		node = node.father
		heap.Push(&h, &heapNode{node: node, deep: i})
	}
	for i := int32(1); i < int32(len(runes)); i++ {
		node, deep := t.findByRunes(runes[i:])
		heap.Push(&h, &heapNode{node: node, deep: deep})
	}
	sz := int32(0)
	querys := make([]*heapNode, 0)
	fmt.Println(h.Len())
	for h.Len() > 0 { // max-heap
		node := heap.Pop(&h)
		querys = append(querys, node.(*heapNode))
		sz += node.(*heapNode).node.size
		if sz > 10 {
			break
		}
	}

	//go query.insertNode()
	for _, node := range querys { // 按deep排序的node数组
		query.getRelatedSearch(node.node, int32(len(runes))-node.deep)
	}

	strings := make([]string, 0)
	for _, node := range query.heap {
		strings = append(strings, *getPrefix(node.node))
	}

	return strings
}

func (q *Query) getRelatedSearch(node *Node, deep int32) {
	for _, v := range node.child { // 是否考虑作为整句的结尾可以获取相对更大的一个优先级
		// 不要残句。。
		if v.count != 0 {
			if q.heap.Len() < 10 {
				heap.Push(&q.heap, &heapNode{node: v, deep: deep})
			} else {
				node := &heapNode{node: v, deep: deep}
				if !q.heap[0].compare(node) {
					heap.Push(&q.heap, node) // 先进后出
					heap.Pop(&q.heap)
				}
			}
		}

		// 维护最大值优化
		maxNode := &heapNode{node: &Node{count: v.max}, deep: deep}
		if q.heap.Len() < 10 || !q.heap[0].compare(maxNode) {
			q.getRelatedSearch(v, deep)
		}
	}
}
