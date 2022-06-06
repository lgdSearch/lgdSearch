package trie

import (
	"bufio"
	"bytes"
	"compress/flate"
	"encoding/gob"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

/*
go test -bench=BenchmarkWrite1 -benchmem -benchtime=1x
BenchmarkWrite1-8              1         129803500 ns/op            4840 B/op          6 allocs/op
go test -bench=BenchmarkWrite2 -benchmem -benchtime=1x
BenchmarkWrite2-8              1         891270700 ns/op        376424000 B/op       180 allocs/op
go test -bench=BenchmarkWrite3 -benchmem -benchtime=1x
BenchmarkWrite3-8              1        3652478100 ns/op        440159312 B/op       214 allocs/op

go test -bench=BenchmarkLoad1 -benchmem -benchtime=1x
BenchmarkLoad1-8               1        7509717100 ns/op        4408894248 B/op 74294868 allocs/op
go test -bench=BenchmarkLoad2 -benchmem -benchtime=1x
BenchmarkLoad2-8               1        8000719200 ns/op        4579538040 B/op 74299381 allocs/op
go test -bench=BenchmarkLoad3 -benchmem -benchtime=1x
BenchmarkLoad3-8               1        8019907800 ns/op        4679763176 B/op 74316042 allocs/op

go test -bench=BenchmarkInsert -benchmem -benchtime=1x
BenchmarkInsert-8              1        6844154500 ns/op        4339703808 B/op 71422123 allocs/op

go test -bench=BenchmarkLoad1 -benchmem -benchtime=1x # with notInsert
BenchmarkLoad1-8               1         228871400 ns/op        70222648 B/op    2876742 allocs/op
go test -bench=BenchmarkLoad2 -benchmem -benchtime=1x # with notInsert
BenchmarkLoad2-8               1         195824000 ns/op        239814776 B/op   2877103 allocs/op
go test -bench=BenchmarkLoad3 -benchmem -benchtime=1x # with notInsert
BenchmarkLoad3-8               1         622688500 ns/op        340043000 B/op   2893798 allocs/op
事实证明序列化没有意义
*/

func write1(filepath string, str []string, counts []int32) {
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
	for _, val := range str {
		_, err := writer.WriteString(val)
		if err != nil {
			return
		}
	}
}

func write2(filepath string, str []string, counts []int32) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(str)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filepath, buffer.Bytes(), 0600)
	if err != nil {
		panic(err)
	}
}

func write3(filepath string, str []string, counts []int32) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(str)
	if err != nil {
		panic(err)
	}

	compressData := Compression(buffer.Bytes())
	err = ioutil.WriteFile(filepath, compressData, 0600)
	if err != nil {
		panic(err)
	}
}

func Compression(data []byte) []byte {
	buf := new(bytes.Buffer)
	write, err := flate.NewWriter(buf, flate.DefaultCompression)
	defer write.Close()

	if err != nil {
		panic(err)
	}

	write.Write(data)
	write.Flush()
	//log.Println("原大小：", len(data), "压缩后大小：", buf.Len(), "压缩率：", fmt.Sprintf("%.2f", float32(buf.Len())*100/float32(len(data))), "%")
	return buf.Bytes()
}

func BenchmarkWrite1(b *testing.B) {
	// run the fib(5) function b.N times
	filepath := "../data/trieData.txt"
	trie := Load(filepath)
	str, counts := serialize(trie.root)
	log.Println(len(str))

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		write1("../data/trieData.txt", str, counts)
	}
}

func BenchmarkWrite2(b *testing.B) {
	filepath := "../data/trieData.txt"
	trie := Load(filepath)
	str, counts := serialize(trie.root)
	log.Println(len(str))

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		write2("../data/trieData2.txt", str, counts)
	}
}

func BenchmarkWrite3(b *testing.B) {
	filepath := "../data/trieData.txt"
	trie := Load(filepath)
	str, counts := serialize(trie.root)
	log.Println(len(str))

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		write3("../data/trieData3.txt", str, counts)
	}
}

func Load1(filepath string) {
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
	//trie := NewTrie()
	cnt := 0
	for {
		cnt += 1
		_, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		//trie.insertRunesWithCount([]rune(str), 1)
	}
	log.Println(cnt)
}

func Load2(filepath string) {
	raw, err := ioutil.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			//忽略
			return
		}
		panic(err)
	}

	data := make([]string, 0)
	buffer := bytes.NewBuffer(raw)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(&data)
	//trie := NewTrie()
	log.Println(len(data))
	//for _, val := range data {
	//	trie.insertRunesWithCount([]rune(val), 1)
	//}
	if err != nil {
		panic(err)
	}
}

func Load3(filepath string) {
	raw, err := ioutil.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			//忽略
			return
		}
		panic(err)
	}

	//解压
	decoData := Decompression(raw)

	data := make([]string, 0)
	buffer := bytes.NewBuffer(decoData)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(&data)
	//trie := NewTrie()
	log.Println(len(data))
	//for _, val := range data {
	//	trie.insertRunesWithCount([]rune(val), 1)
	//}
	if err != nil {
		panic(err)
	}
}

//Decompression 解压缩数据
func Decompression(data []byte) []byte {
	return DecompressionBuffer(data).Bytes()
}

func DecompressionBuffer(data []byte) *bytes.Buffer {
	buf := new(bytes.Buffer)
	read := flate.NewReader(bytes.NewReader(data))
	defer read.Close()

	buf.ReadFrom(read)
	return buf
}

func BenchmarkLoad1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Load1("../data/trieData.txt")
	}
}

func BenchmarkLoad2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Load2("../data/trieData2.txt")
	}
}

func BenchmarkLoad3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Load3("../data/trieData3.txt")
	}
}

func Insert(trie *Trie, str []string) {
	for _, val := range str {
		go trie.insertRunesWithCount([]rune(val), 1)
	}
}

func BenchmarkInsert(b *testing.B) {
	raw, err := ioutil.ReadFile("../data/trieData3.txt")
	if err != nil {
		if os.IsNotExist(err) {
			//忽略
			return
		}
		panic(err)
	}

	//解压
	decoData := Decompression(raw)

	data := make([]string, 0)
	buffer := bytes.NewBuffer(decoData)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(&data)
	trie := NewTrie()
	log.Println(len(data))
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		Insert(trie, data)
	}
	if err != nil {
		panic(err)
	}
}

func getPrefixDFS(node *Node) *string {
	str := ""
	for node.father != nil {
		defer func(str *string, node *Node) {
			*str += string(node.data)
		}(&str, node)
		node = node.father
	}
	return &str
}

func getPrefixNODFS(node *Node) *string {
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

// BenchmarkGetPrefixDFS-8           969994              1216 ns/op             568 B/op         30 allocs/op
func BenchmarkGetPrefixDFS(b *testing.B) {
	root := newNode()
	for i := 0; i < 10; i++ { // 假设平均长度为 10
		node := newNode()
		node.data = '我'
		node.father = root
		root.child[node.data] = node

		root = node
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		getPrefixDFS(root)
	}
}

// BenchmarkGetPrefixNODFS-8        4501153               273.4 ns/op           168 B/op          5 allocs/op
func BenchmarkGetPrefixNODFS(b *testing.B) {
	root := newNode()
	for i := 0; i < 10; i++ { // 假设平均长度为 10
		node := newNode()
		node.data = '我'
		node.father = root
		root.child[node.data] = node

		root = node
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		getPrefixNODFS(root)
	}
}
