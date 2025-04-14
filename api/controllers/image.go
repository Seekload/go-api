package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// ImageRequest 定义请求结构
type ImageRequest struct {
	Files       []string `json:"files"` // 本地文件路径
	Prompt      string   `json:"prompt"`
	Size        string   `json:"size"`
	CallbackURL string   `json:"callBackUrl"`
}

// GenerateImage 处理图片生成请求
func GenerateImage(c *gin.Context) {
	var req ImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 获取项目根目录
	rootDir, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取项目根目录失败"})
		return
	}

	// 验证文件是否存在并收集文件URL
	var fileURLs []string
	for _, file := range req.Files {
		// 构建完整的文件路径
		fullPath := filepath.Join(rootDir, file)
		log.Println("Checking file:", fullPath)

		// 检查文件是否存在
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "文件不存在: " + file})
			return
		}

		// 将本地路径转换为 file:// URL
		fileURLs = append(fileURLs, "file://"+fullPath)
	}

	// 构建请求体
	requestBody := map[string]interface{}{
		"filesUrl":    fileURLs,
		"prompt":      req.Prompt,
		"size":        req.Size,
		"callBackUrl": req.CallbackURL,
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "构建请求体失败"})
		return
	}

	// 创建 HTTP 请求
	client := &http.Client{}
	request, err := http.NewRequest("POST", "https://kieai.erweima.ai/api/v1/gpt4o-image/generate", bytes.NewBuffer(payload))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建请求失败"})
		return
	}

	// 设置请求头
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer 1c7da3dd8bc930d25a55733cdaa24e27")

	// 发送请求
	response, err := client.Do(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送请求失败"})
		return
	}
	defer response.Body.Close()

	// 读取响应
	body, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取响应失败"})
		return
	}

	// 返回响应
	c.Data(response.StatusCode, "application/json", body)
}
