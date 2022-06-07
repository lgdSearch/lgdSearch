# lgdSearch

**字节跳动青训营项目**

`lgdSearch` 一个golang实现的全文检索引擎，支持持久化，相关搜索，搜索文档和图片。

## 文档

+ [API文档](./docs/swagger.yaml) 使用 [swagger](https://editor.swagger.io/) 查看
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
    |   |-- db                 数据库操作的封装
    |   |   |-- badgerStorage
    |   |   |-- boltStorage
    |   |-- extractclaims      jwtMapClaims解析
    |   |-- httprequest        http请求构建（测试用）
    |   |-- logger             日志持久化
    |   |-- models             数据库表结构
    |   |-- pagination
    |   |-- trie               各类树的实现
    |   |-- utils
    |   |   |-- colf           colf序列化工具
    |   |       |-- doc        doc序列化
    |   |       |-- keyIds     keyIds序列化
    |   |-- vgg                图片特征值提取
    |   |-- weberror           接口返回错误用时使用的结构体
    |-- router                 gin接口路由
```
