# lgdSearch

**字节跳动青训营项目**

`lgdSearch` 一个golang实现的全文检索引擎，支持持久化，相关搜索，搜索文档和图片。

## 文档

+ [API文档](./docs/swagger.yaml) 使用 [swagger](https://editor.swagger.io/) 查看
+ [搜文档|搜图片](./docs/search.md)
+ [相关搜索](./docs/related_search.md)
+ [持久化](./docs/storage.md)
+ [以图搜图](./docs/imageSearch.md)

## 技术栈

+ Trie树
+ AC自动机
+ 快速排序法
+ 倒排索引
+ 文件分片
+ golang-jieba分词
+ boltdb
+ badgerdb
+ colfer序列化
+ grpc
+ ResNet50
+ diskcache
+ milvus

## 目录结构
```
|-- lgdSearch
    |-- controller             http接口
    |-- docs                   文档
    |-- handler                CRUD操作
    |-- logs                   持久化日志
    |-- middleware             gin中间件
    |-- payloads               接口req与resp结构体
    |-- pkg
    |   |-- data               存储各类数据
    |   |   |-- badger_doc_0.db存储文档数据
    |   |   |-- ....
    |   |   |-- picture        存储缩略图信息
    |   |   |-- bolt_keyIds.db 存储倒排索引
    |   |   |-- dataIndex.txt  持久化数据库中文档数
    |   |   |-- dictionary.txt jieba分词文档
    |   |   |-- HotSearch.txt  持久化热搜
    |   |   |-- trieData.txt   持久化Trie树
    |   |-- db                 数据库操作的封装
    |   |   |-- badgerStorage  
    |   |   |-- boltStorage
    |   |-- extractclaims      jwtMapClaims解析
    |   |-- httprequest        http请求构建（测试用）
    |   |-- logger             日志持久化
    |   |-- models             数据库表结构
    |   |-- pagination         分页功能
    |   |-- trie               各类树和数据结构的实现
    |   |-- utils  
    |   |   |-- colf           colf序列化工具库
    |   |       |-- doc        doc序列化
    |   |       |-- keyIds     keyIds序列化
    |   |-- vgg                以图搜图
    |   |-- weberror           接口返回错误用时使用的结构体
    |-- router                 gin接口路由
```

## 在线体验
在 2022/9 月之前可以通过 [这里](http://121.196.207.80:8081) 在线体验本项目
