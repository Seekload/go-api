package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// ImageRequest 定义请求结构
type ImageRequest struct {
	ImageUrls   []string `json:"imageUrls" binding:"required"`   // 线上图片URL列表
	Prompt      string   `json:"prompt" binding:"required"`      // 提示词
	Size        string   `json:"size" binding:"required"`        // 图片尺寸
	CallbackURL string   `json:"callBackUrl" binding:"required"` // 回调地址
}

// GhibliResponse 定义吉卜力图片生成的响应结构
type GhibliResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	TaskID   string `json:"taskId,omitempty"`
	ImageURL string `json:"imageUrl,omitempty"`
}

// RemoveBackgroundRequest 定义背景移除请求结构
type RemoveBackgroundRequest struct {
	ImageURL string `json:"imageUrl,omitempty"` // 可选：图片URL
}

// RemoveBackgroundResponse 定义背景移除响应结构
type RemoveBackgroundResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	ImageURL string `json:"imageUrl,omitempty"`
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

// isValidImageType 验证图片类型
func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
	}

	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}

// RemoveBackground 处理背景移除请求
func RemoveBackground(c *gin.Context) {
	// 检查是否有上传的文件
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		// 如果没有文件，检查是否有图片URL
		var req RemoveBackgroundRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.ImageURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "请上传图片文件或提供图片URL",
			})
			return
		}

		// 使用URL处理背景移除
		resultURL, err := removeBackgroundFromURL(req.ImageURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "背景移除失败: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"message":  "背景移除成功",
			"imageUrl": resultURL,
		})
		return
	}
	defer file.Close()

	// 验证文件类型
	if !isValidImageType(header.Header.Get("Content-Type")) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "只支持 JPG、PNG、GIF 格式的图片",
		})
		return
	}

	// 调用Photoroom API移除背景
	resultURL, err := removeBackgroundFromFile(file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "背景移除失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "背景移除成功",
		"imageUrl": resultURL,
	})
}

// removeBackgroundFromFile 从上传的文件移除背景
func removeBackgroundFromFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	// 获取Photoroom API密钥
	apiKey := os.Getenv("PHOTOROOM_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("未配置Photoroom API密钥")
	}

	// 创建multipart请求
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 添加图片文件
	part, err := writer.CreateFormFile("image_file", header.Filename)
	if err != nil {
		return "", fmt.Errorf("创建表单文件失败: %v", err)
	}

	// 重置文件指针并复制文件内容
	file.Seek(0, 0)
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("复制文件内容失败: %v", err)
	}

	writer.Close()

	// 创建HTTP请求
	req, err := http.NewRequest("POST", "https://sdk.photoroom.com/v1/segment", &requestBody)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-API-Key", apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(bodyBytes))
	}

	// 检查响应是否为图片
	contentType := resp.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("期望图片响应，但收到: %s", string(bodyBytes))
	}

	// 读取处理后的图片
	processedImage, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取处理后的图片失败: %v", err)
	}

	// 上传处理后的图片到Vercel Blob Storage
	imageURL, err := uploadProcessedImageToBlob(processedImage, header.Filename)
	if err != nil {
		return "", fmt.Errorf("上传处理后的图片失败: %v", err)
	}

	return imageURL, nil
}

// removeBackgroundFromURL 从URL移除背景
func removeBackgroundFromURL(imageURL string) (string, error) {
	// 下载图片
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("下载图片失败: %v", err)
	}
	defer resp.Body.Close()

	// 获取Photoroom API密钥
	apiKey := os.Getenv("PHOTOROOM_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("未配置Photoroom API密钥")
	}

	// 创建multipart请求
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 添加图片URL
	err = writer.WriteField("image_url", imageURL)
	if err != nil {
		return "", fmt.Errorf("添加图片URL失败: %v", err)
	}

	writer.Close()

	// 创建HTTP请求
	req, err := http.NewRequest("POST", "https://sdk.photoroom.com/v1/segment", &requestBody)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-API-Key", apiKey)

	// 发送请求
	client := &http.Client{}
	photoroomResp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer photoroomResp.Body.Close()

	// 检查响应状态
	if photoroomResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(photoroomResp.Body)
		return "", fmt.Errorf("API请求失败，状态码: %d, 响应: %s", photoroomResp.StatusCode, string(bodyBytes))
	}

	// 读取处理后的图片
	processedImage, err := io.ReadAll(photoroomResp.Body)
	if err != nil {
		return "", fmt.Errorf("读取处理后的图片失败: %v", err)
	}

	// 上传处理后的图片到Vercel Blob Storage
	filename := "removed_bg_" + fmt.Sprintf("%d", time.Now().Unix()) + ".png"
	imageURL, err = uploadProcessedImageToBlob(processedImage, filename)
	if err != nil {
		return "", fmt.Errorf("上传处理后的图片失败: %v", err)
	}

	return imageURL, nil
}

// uploadProcessedImageToBlob 上传处理后的图片到Vercel Blob Storage
func uploadProcessedImageToBlob(imageData []byte, filename string) (string, error) {
	// 获取Vercel Blob Token
	blobToken := os.Getenv("BLOB_READ_WRITE_TOKEN")
	if blobToken == "" {
		return "", fmt.Errorf("未配置Vercel Blob存储令牌")
	}

	// 生成唯一的文件名
	uniqueFilename := fmt.Sprintf("bg_removed_%d_%s", time.Now().Unix(), filename)

	// 准备上传数据
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 添加文件数据
	part, err := writer.CreateFormFile("file", uniqueFilename)
	if err != nil {
		return "", fmt.Errorf("创建表单文件失败: %v", err)
	}

	_, err = part.Write(imageData)
	if err != nil {
		return "", fmt.Errorf("写入文件数据失败: %v", err)
	}

	writer.Close()

	// 构建上传URL - 这里需要你的实际Blob Store URL
	// 格式通常是: https://[store-id].public.blob.vercel-storage.com/
	// 在实际使用中，你需要从环境变量或配置中获取正确的URL
	storeURL := os.Getenv("VERCEL_BLOB_STORE_URL")
	if storeURL == "" {
		return "", fmt.Errorf("未配置Vercel Blob存储URL")
	}

	// 创建HTTP请求
	req, err := http.NewRequest("PUT", storeURL+"/"+uniqueFilename, &requestBody)
	if err != nil {
		return "", fmt.Errorf("创建上传请求失败: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+blobToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("上传请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("上传失败，状态码: %d, 响应: %s", resp.StatusCode, string(bodyBytes))
	}

	// 解析响应获取URL
	var uploadResponse struct {
		URL string `json:"url"`
	}

	var responseBody []byte
	responseBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	err = json.Unmarshal(responseBody, &uploadResponse)
	if err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if uploadResponse.URL == "" {
		return "", fmt.Errorf("上传响应中未包含URL")
	}

	return uploadResponse.URL, nil
}
