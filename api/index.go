package handler

import (
	"go-api/api/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// 创建 Gin 引擎
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// 注册路由
	routes.SetupRoutes(router)

	// 处理请求
	router.ServeHTTP(w, r)
}
