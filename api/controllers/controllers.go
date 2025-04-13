package controllers

import (
	"github.com/gin-gonic/gin"
)

func Hello(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, "<h1>Hello from www.zzfly.net!</h1>")
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
