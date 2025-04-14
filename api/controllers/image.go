package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ImageRequest 定义请求结构
type ImageRequest struct {
	ImageUrls   []string `json:"imageUrls" binding:"required"`   // 线上图片URL列表
	Prompt      string   `json:"prompt" binding:"required"`      // 提示词
	Size        string   `json:"size" binding:"required"`        // 图片尺寸
	CallbackURL string   `json:"callBackUrl" binding:"required"` // 回调地址
}

// GenerateImage 处理图片生成请求
func GenerateImage(c *gin.Context) {
	var req ImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 验证图片URL列表
	if len(req.ImageUrls) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "至少需要一张图片"})
		return
	}

	// 验证回调地址
	if req.CallbackURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "回调地址不能为空"})
		return
	}

	// 构建请求体
	requestBody := map[string]interface{}{
		"filesUrl":    req.ImageUrls,
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

// GetTaskInfo 获取任务信息
func GetTaskInfo(c *gin.Context) {
	// 获取任务ID
	taskId := c.Query("taskId")
	if taskId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "任务ID不能为空"})
		return
	}

	// 创建 HTTP 请求
	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://kieai.erweima.ai/api/v1/gpt4o-image/record-info", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建请求失败"})
		return
	}

	// 设置请求头
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer 1c7da3dd8bc930d25a55733cdaa24e27")

	// 添加查询参数
	q := request.URL.Query()
	q.Add("taskId", taskId)
	request.URL.RawQuery = q.Encode()

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
