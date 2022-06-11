package trie

import (
	"fmt"
	"log"
	"testing"
)

func TestTrie(t *testing.T) {
	trie := Load("../data/trieData.txt")
	log.Println("trie 初始化完成")

	str := "你好我是你爷爷"
	trie.InsertString(str)

	str = "回复"
	result := trie.Search([]rune(str))
	fmt.Println(len(result))
	for _, val := range result {
		fmt.Print(val)
	}

	trie.FlushIndex("../data/trieData.txt")
	log.Println("trie 写入成功")
}
