package routes

import (
	"go-api/api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// 添加中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 定义路由组
	api := router.Group("/api")
	{
		// 示例路由
		api.GET("/hello", controllers.Hello)
		api.GET("/ping", controllers.Ping)

		// 可以在这里添加更多路由组
		// v1 := api.Group("/v1")
		// {
		//     v1.GET("/users", controllers.GetUsers)
		//     v1.POST("/users", controllers.CreateUser)
		// }
	}
}
