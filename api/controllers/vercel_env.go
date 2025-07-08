package controllers

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// VercelEnvDebug 检查Vercel特定的环境变量
func VercelEnvDebug(c *gin.Context) {
	// 获取所有环境变量
	allEnvs := os.Environ()

	// 过滤Vercel相关的环境变量
	vercelEnvs := make(map[string]string)
	r2Envs := make(map[string]string)

	for _, env := range allEnvs {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]

			if strings.Contains(strings.ToUpper(key), "VERCEL") {
				vercelEnvs[key] = maskValue(value)
			}

			if strings.HasPrefix(key, "R2_") || key == "PHOTOROOM_API_KEY" {
				r2Envs[key] = maskValue(value)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Vercel环境变量详细检查",
		"vercel_env_count": len(vercelEnvs),
		"vercel_envs":      vercelEnvs,
		"r2_envs":          r2Envs,
		"total_env_count":  len(allEnvs),
	})
}

// maskValue 隐藏敏感值
func maskValue(value string) string {
	if len(value) == 0 {
		return "空值"
	}
	if len(value) <= 4 {
		return "***"
	}
	if len(value) <= 10 {
		return value[:2] + "***"
	}
	return value[:4] + "***" + value[len(value)-4:]
}
