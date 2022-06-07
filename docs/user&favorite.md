# 用户与收藏夹

使用gorm操作mysql储存数据，gin-jwt生成鉴权中间件

具体接口见[api文档](swagger.json)

## 用户

表结构
```
type User struct {
	gorm.Model
	Username  string
	Nickname  string
	Password  string
	Favorites []Favorite //收藏夹
}
```
登录与登出功能使用gin-jwt生成接口
还有注册、删除账户、修改昵称、获取个人信息功能

## 收藏夹
表结构
```
//收藏夹
type Favorite struct {
	gorm.Model
	UserId 	    uint
	Name        string
	Docs 	    []Doc
}

//收藏
type Doc struct {
	gorm.Model
	FavoriteId  uint
	DocIndex    uint    //搜索接口返回的文档号
	Summary     string  //文档前一部分字符
}
```

有创建收藏夹、收藏夹重命名、获取收藏夹内容等功能
获取收藏夹列表与收藏夹内容支持分页