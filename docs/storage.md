# 持久化

## 倒排索引

倒排索引通过 boltdb 存储, 主要是看重了 boltdb 随机读的高效率。
boltdb 是基于 B+ 树结构实现的存储结构，所以查找的效率很高，但是修改的效率非常低下。

以下为倒排索引数组中存储的对象结构，使用 colfer 序列化

[点击查看源码](../pkg/utils/colf/keyIds/keyIds.go)
```go
package keyIds

// boltdb 中的 Ids 对象
type StorageIds struct {
	StorageIds []*StorageId
}

// StorageId boltdb中的Id存储对象
type StorageId struct {
	Id uint32

	Score float32 // 这个分词在文档(id)中的分数
}
```

## 文档数据

文档数据通过 badgerdb 存储，
主要是看重了 badgerdb 对于读和写的优秀平衡。 

badgerdb 是leveldb的改良版本，
写入速度媲美 leveldb 的同时改进了 leveldb 原本过慢的随机读效率。

这里不使用boltdb的原因是文档相对于倒排索引对随机读的速度的要求没那么高，
同时 boltdb 的写入效率过低会导致写文档非常缓慢

以下是文档对象的结构，对于文档的ID，
因为可以通过 k-v 键值对的 key 记录，
所以在 value 中不记录

使用 colfer 序列化

[点击查看源码](../pkg/utils/colf/doc/docStorage.go)
```go
package doc

// StorageIndexDoc 文档对象
type StorageIndexDoc struct {
	Text string

	Url string
}
```
