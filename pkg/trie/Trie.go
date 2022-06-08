// Copyright (c) 2022, wan-she li. All rights reserved

package trie

import (
	"container/heap"
	"sync"
)

const TrieDataPath = "trieData.txt"

var Tree *Trie

func InitTrie(filepath string) {
	Tree = Load(filepath)
	go Tree.automaticFlush(filepath) // 刷新到磁盘
}

type pair struct {
	key   int32
	value *Node
}

// Node is a single element within the Trie
type Node struct {
	father *Node // 双向关系
	data   rune
	sons   []pair         // 数量少于9使用这个，否则使用map
	child  map[rune]*Node // 大于8以后申请map，将sons数据放入map，sons清空

	size  int32 // 统计子树完整串数量
	count int32 // 累计查找次数(用于排序)，这个是以这个词为结尾的句子个数 : int32 = rune
	max   int32 // 维护子树中最大 count, 部分节点可能 count = 0, 但是 max 有值，这表明它是一个非结尾，但是它下面有结尾
}

// get a new Node
func newNode() *Node {
	return &Node{}
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

	stk := make([]*Node, 0)
	flag := false
	for _, val := range runes {
		var nxt *Node
		var ok = false
		if now.child != nil { // map
			nxt, ok = now.child[val]
			if !ok {
				nxt = newNode()
				nxt.data = val
				now.child[val] = nxt
				nxt.father = now
				t.size += 1
			}
		} else { // sons
			if now.sons != nil {
				for _, pair := range now.sons {
					if pair.key == val {
						nxt, ok = pair.value, true
					}
				}
			}
			if !ok {
				nxt = newNode()
				nxt.data = val
				nxt.father = now
				t.size += 1
				if now.sons == nil {
					now.sons = make([]pair, 1)
					now.sons[0] = pair{value: nxt, key: val}
				} else if len(now.sons) < 9 {
					now.sons = append(now.sons, pair{value: nxt, key: val})
				} else { // > 8 转map
					now.child = make(map[rune]*Node)
					now.child[val] = nxt
					for _, pair := range now.sons {
						now.child[pair.key] = pair.value
					}
					now.sons = nil // gc
				}
			}
		}
		stk = append(stk, now)
		now = nxt
	}
	if now.count == 0 { // 第一次成为完整串, 贡献一个 size
		now.size += 1
		flag = true
	}
	now.count += count // 此句出现count次

	// 维护每个节点的 max
	if flag {
		for i := len(runes) - 2; i >= 0; i-- {
			now.size += 1
			if stk[i].max < stk[i+1].max {
				stk[i].max = stk[i+1].max
			}
		}
	} else {
		for i := len(runes) - 2; i >= 0; i-- {
			if stk[i].max < stk[i+1].max {
				stk[i].max = stk[i+1].max
			} else {
				break
			}
		}
	}
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
	now := t.root
	deep := int32(0)
	var nxt *Node
	var ok = false
	for _, val := range runes {
		if now.child != nil {
			nxt, ok = now.child[val]
		} else {
			ok = false
			for _, pair := range now.sons {
				if pair.key == val {
					nxt = pair.value
					ok = true
				}
			}
		}
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
	heap Heap
	sync.WaitGroup
}

// 获取前缀
func getPrefix(node *Node) *string {
	str := make([]rune, 0)
	for node.father != nil {
		str = append(str, node.data)
		node = node.father
	}
	for i, j := 0, len(str)-1; i < j; i, j = i+1, j-1 {
		str[i], str[j] = str[j], str[i]
	}
	result := string(str)
	return &result
}

// Search 查找相关搜索词条, 默认返回不大于十条
func (t *Trie) Search(runes []rune) []string {
	query := &Query{heap: Heap{}}
	node, deep := t.findByRunes(runes)

	h := MaxHeap{}
	if node != nil { // 判空很重要
		heap.Push(&h, &heapNode{node: node, deep: deep})
	}
	sumSz := node.size
	data := node.data
	for i := deep - 1; i > 0; i-- {
		node = node.father
		if node == nil { // nil 判空
			break
		}
		if node.size > sumSz {
			heap.Push(&h, &heapNode{node: node, deep: i, size: node.size - sumSz, son: data})
		}
		sumSz = node.size
		data = node.data
	}
	for i := int32(1); i < int32(len(runes)); i++ {
		node, deep := t.findByRunes(runes[i:])
		if node != nil {
			heap.Push(&h, &heapNode{node: node, deep: deep})
		}
	}
	sz := int32(0)
	querys := make([]*heapNode, 0)
	for h.Len() > 0 { // max-heap
		node := heap.Pop(&h)
		querys = append(querys, node.(*heapNode))
		sz += node.(*heapNode).size
		if sz > 10 {
			break
		}
	}

	for _, node := range querys { // 按deep排序的node数组
		query.getRelatedSearch(node.son, node.node, int32(len(runes))-node.deep)
	}

	strings := make([]string, 0)
	for _, node := range query.heap {
		strings = append(strings, *getPrefix(node.node))
	}

	return strings
}

// 遍历以 node 为根的整颗子树, 将完整词条的节点插入堆中
func (q *Query) getRelatedSearch(son rune, node *Node, deep int32) {
	if node.child != nil {
		for k, v := range node.child {
			if son != 0 && k == son { // 跳过儿子
				continue
			}
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
				q.getRelatedSearch(0, v, deep)
			}
		}
	} else {
		for _, pair := range node.sons {
			if son != 0 && pair.value.data == son { // 跳过儿子
				continue
			}
			if pair.value.count != 0 {
				if q.heap.Len() < 10 {
					heap.Push(&q.heap, &heapNode{node: pair.value, deep: deep})
				} else {
					node := &heapNode{node: pair.value, deep: deep}
					if !q.heap[0].compare(node) {
						heap.Push(&q.heap, node) // 先进后出
						heap.Pop(&q.heap)
					}
				}
			}

			//维护最大值优化
			maxNode := &heapNode{node: &Node{count: pair.value.max}, deep: deep}
			if q.heap.Len() < 10 || !q.heap[0].compare(maxNode) {
				q.getRelatedSearch(0, pair.value, deep)
			}
		}
	}
}
