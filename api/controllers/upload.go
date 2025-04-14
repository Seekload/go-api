package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadImage 处理图片上传请求
func UploadImage(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的文件"})
		return
	}

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "upload-*.tmp")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建临时文件失败"})
		return
	}
	defer os.Remove(tempFile.Name())

	// 保存上传的文件到临时文件
	if err := c.SaveUploadedFile(file, tempFile.Name()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 读取文件内容
	fileData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取文件失败"})
		return
	}

	// 生成唯一的文件名
	filename := time.Now().Format("20060102150405") + "-" + file.Filename

	// 上传到 Vercel Blob
	blobURL, err := uploadToVercelBlob(fileData, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传到 Blob Storage 失败"})
		return
	}

	// 返回上传结果
	c.JSON(http.StatusOK, gin.H{
		"url": blobURL,
	})
}

// uploadToVercelBlob 上传文件到 Vercel Blob Storage
func uploadToVercelBlob(data []byte, filename string) (string, error) {
	// 创建上传请求
	client := &http.Client{}
	req, err := http.NewRequest("PUT", "https://blob.vercel-storage.com/"+filename, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("x-vercel-blob-filename", filename)
	req.Header.Set("x-vercel-blob-content-type", "image/jpeg") // 根据实际文件类型设置
	req.Header.Set("Authorization", "Bearer vercel_blob_rw_voQSUpJHWEa9uhr2_0DIzwctA5pysVuWIRy3HmWKc8YJeuL")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析响应
	var result struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.URL, nil
}
