package trie

import (
	"bufio"
	"io"
	"os"
	"runtime"
	"time"
)

var strings []string
var counts []int32

// 用于遍历子树
func (n *Node) foreach(runes []rune, deep int) {
	if n.count != 0 {
		strings = append(strings, string(runes))
		counts = append(counts, n.count)
	}
	for _, v := range n.child {
		if len(runes) <= deep {
			runes = append(runes, '我')
		}
		runes[deep] = v.data
		v.foreach(runes, deep+1)
	}
}

// serialize
// node 需要传入根节点
// []string, []int32 返回 trie 树中所有的完整词条以及对应的count值
func serialize(node *Node) ([]string, []int32) {
	strings, counts = make([]string, 0), make([]int32, 0)
	if node == nil {
		return strings, counts
	}

	runes := make([]rune, 0)
	node.foreach(runes, 0)

	return strings, counts
}

// Write 序列化写入文件
func Write(trie *Trie, filepath string) {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic("write trie error: can't open file!")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic("write trie error: file close fail")
		}
	}(file)

	writer := bufio.NewWriter(file)
	str, counts := serialize(trie.root)
	for pos, val := range str {
		_, err := writer.WriteString(val + string(counts[pos]))
		if err != nil {
			return
		}
	}
}

// Load 从文件中加载 trie 树
func Load(filepath string) *Trie {
	file, err := os.Open(filepath)
	if err != nil {
		panic("Trie load error: can't find file!")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic("trie load error: file close fail")
		}
	}(file)

	reader := bufio.NewReader(file)
	trie := NewTrie()
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		val, _, err := reader.ReadRune()                      // _ 是val作为rune的实际字节长度
		go trie.insertRunesWithCount([]rune(str), int32(val)) // 强制类型转化的开销大吗，还是说直接将 []rune 存进去，读出来
	}

	return trie
}

// 自动保存索引，60秒钟检测一次
func (t *Trie) automaticFlush(filepath string) {
	ticker := time.NewTicker(time.Second * 60)
	size := 0

	for {
		<-ticker.C
		//检查数据是否有变动
		if size != t.size {
			size = t.size
			t.FlushIndex(filepath)
		}
		//定时GC
		runtime.GC()
	}

}

// FlushIndex 刷新缓存到磁盘
func (t *Trie) FlushIndex(filepath string) {
	t.Lock()
	defer t.Unlock()

	Write(t, filepath)
}
