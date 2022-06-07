package router

import (
	"github.com/gin-gonic/gin"
	"lgdSearch/controller"
	"lgdSearch/middleware"
	"net/http"
)

func Init() *gin.Engine {
	//设置默认引擎
	engine := gin.Default()
	engine.Use(cors())
	auth := middleware.GetJWTMiddle()
	engine.POST("/login", auth.LoginHandler)
	engine.PUT("/register", controller.Register)

	users := engine.Group("/users")
	users.Use(auth.MiddlewareFunc())
	{
		users.PUT("/logout", auth.LogoutHandler)
		users.DELETE("", controller.DeleteAccount)
		users.PATCH("/nickname", controller.UpdateNickname)
		users.GET("/profile", controller.GetProfile)
		favorites := users.Group("/favorites")
		{
			favorites.PUT("", controller.AddFavorite)
			favorites.PATCH("/:fav_id/name", controller.UpdateFavoriteName)
			favorites.DELETE("/:fav_id", controller.DeleteFavorite)
			favorites.GET("/:fav_id", controller.GetFavorite)
			favorites.GET("", controller.GetFavorites)
			favorites.PUT("/:fav_id/docs", controller.AddDoc)
			favorites.DELETE("/:fav_id/docs/:doc_id", controller.DeleteDoc)
			favorites.GET("/:fav_id/docs", controller.GetDocs)
		}
	}
	engine.POST("/query", controller.Query)                // 查询
	engine.POST("/query/picture", controller.QueryPicture) // 查询图片
	engine.GET("/query/hotSearch", controller.HotSearch)   // 热搜
	engine.GET("/book/:text", controller.GetRelatedSearch) // 相关搜索
	return engine
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Header("Access-Control-Allow-Headers", "Content-Type,X-CSRF-Token, Authorization, Token,Access-Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,PUT,DELETE,PATCH")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
