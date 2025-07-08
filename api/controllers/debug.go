package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// DebugEnv 调试环境变量（仅用于测试，生产环境请删除）
func DebugEnv(c *gin.Context) {
	// 获取原始环境变量
	rawEnvVars := map[string]interface{}{
		"R2_ACCOUNT_ID":        checkEnvVar("R2_ACCOUNT_ID"),
		"R2_ACCESS_KEY_ID":     checkEnvVar("R2_ACCESS_KEY_ID"),
		"R2_SECRET_ACCESS_KEY": checkEnvVar("R2_SECRET_ACCESS_KEY"),
		"R2_BUCKET_NAME":       checkEnvVar("R2_BUCKET_NAME"),
		"R2_DEV_DOMAIN":        checkEnvVar("R2_DEV_DOMAIN"),
		"PHOTOROOM_API_KEY":    checkEnvVar("PHOTOROOM_API_KEY"),
	}

	// 测试R2客户端创建
	r2Client, r2Error := NewR2Client()
	r2Status := map[string]interface{}{
		"success": r2Client != nil,
		"error":   "",
	}
	if r2Error != nil {
		r2Status["error"] = r2Error.Error()
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "环境变量和R2连接检查",
		"env_vars":  rawEnvVars,
		"r2_client": r2Status,
	})
}

// checkEnvVar 检查环境变量状态
func checkEnvVar(key string) map[string]interface{} {
	value := os.Getenv(key)
	return map[string]interface{}{
		"exists":       value != "",
		"masked_value": maskString(value),
		"length":       len(value),
	}
}

// maskString 隐藏敏感信息
func maskString(s string) string {
	if s == "" {
		return "未设置"
	}
	if len(s) <= 8 {
		return "***已设置***"
	}
	return s[:4] + "***" + s[len(s)-4:]
}
