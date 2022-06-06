# lgdSearch

**字节跳动青训营项目**

`lgdSearch` 一个golang实现的全文检索引擎，支持持久化，相关搜索，搜索文档和图片。

## 文档

+ API 使用 [swagger](https://editor.swagger.io/) 查看
+ [搜文档|搜图片](./docs/search.md)
+ [相关搜索](./docs/related_search.md)
+ [持久化](./docs/storage.md)

## 技术栈

+ Trie树
+ AC自动机
+ 快速排序法
+ 倒排索引
+ 文件分片
+ golang-jieba分词
+ boltdb
+ badgerdb
+ colf序列化
