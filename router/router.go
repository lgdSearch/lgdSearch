package router

import (
	"lgdSearch/controller"
	"lgdSearch/middleware"
	"net/http"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	//设置默认引擎
	engine := gin.Default()
	engine.Use(cors())
	auth := middleware.GetJWTMiddle()
	engine.POST("/login", auth.LoginHandler)
	engine.POST("/logout", auth.LogoutHandler)
	engine.POST("/register", controller.Register)
	user := engine.Group("/user")
	user.Use(auth.MiddlewareFunc())
	{
		engine.POST("/update", controller.UpdateProfile)
	}
	engine.GET("/test/say_hello", controller.SayHello)
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