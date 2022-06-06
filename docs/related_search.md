# 相关搜索

## 相关搜索

概述：相关搜索使用 Trie 树实现，对于每个查询的查询语句，插入 Trie 中
，当查询相关搜索时通过相关性算法获得优先级前十的相关语句

以下是 Trie 树及其节点结构[点击查看源码](../pkg/trie/Trie.go)
```go
package trie

import "sync"

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

// Trie holds elements of the Trie tree.
type Trie struct {
	size int
	root *Node
	sync.Mutex
}
```
### Trie 节点关系
Trie树中的节点之间的关系通过两种结构分类维护
* 当子节点数量小于等于8，通过数组`sons`维护，减小内存开销
* 子节点数大于8后申请一个map`child`，将`sons`中所有节点放入`child`中并将sons回收

在Trie树插入字符串数 500W 条情况下，使用纯 map 结构维护需要 `3.5G` 内存，而改用两种结构混合的模式只要`1.5G`内存，同时效率基本无差异

### 相关性算法
相关性算法通过字符串匹配相关度来比较
例如对于串`你好字节跳动`
对于以下Trie中搜索出来的相关串匹配度为

`你好字节跳动` 6

`你好字节飞舞` 4

`字节跳动是啥` 4

通过深度优先搜索来搜索匹配串，同时过程中用一个小顶堆来维护相关度前10的匹配串，
采用一定的剪枝策略加速搜索。详细内容可看[源码和注释](../pkg/trie/Trie.go)

对数据量 500W 条，节点 2000W 个的情况下搜索可达到毫秒级查询

## 热搜

概述：通过一个链表维护24小时内用户的搜索记录，定时更新搜索次数排前十的数据

以下是数据对象结构[点击查看源码](../pkg/trie/HotSearch.go)
```go
package hotSearch

import (
	"sync"
	"time"
)

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

type HotSearch struct {
	searchMessage  map[string]int // 存热点数据的搜索次数
	searchQueue    *queue         // queue 维护map中的数据
	hotSearchArray []*HotSearchMessage
	sync.Mutex
}
```
