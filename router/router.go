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
		users.DELETE("/logout", auth.LogoutHandler)
		users.DELETE("", controller.DeleteAccount)
		users.PATCH("/nickname", controller.UpdateNickname)
		users.GET("/profile", controller.GetProfile)
		favortes := users.Group("/favorites")
		{
			favortes.PUT("/:doc_id", controller.AddFavorite)
			favortes.DELETE("/:doc_id", controller.DeleteFavorite)
			favortes.GET("", controller.GetFavorites)
		}
	}
	engine.GET("/query/hotSearch", controller.HotSearch) // 热搜
	engine.GET("/test/say_hello", controller.SayHello)
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
