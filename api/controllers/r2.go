package controllers

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// R2Config R2存储配置
type R2Config struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
}

// R2Client R2客户端
type R2Client struct {
	client *s3.Client
	config *R2Config
}

// NewR2Client 创建新的R2客户端
func NewR2Client() (*R2Client, error) {
	// 首先尝试从环境变量获取
	cfg := &R2Config{
		AccountID:       os.Getenv("R2_ACCOUNT_ID"),
		AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
		BucketName:      os.Getenv("R2_BUCKET_NAME"),
	}

	// 再次验证配置
	if cfg.AccountID == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" || cfg.BucketName == "" {
		return nil, fmt.Errorf("R2配置无效")
	}

	// 创建AWS配置
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("创建AWS配置失败: %v", err)
	}

	// 创建S3客户端，指向R2端点
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID))
	})

	return &R2Client{
		client: client,
		config: cfg,
	}, nil
}

// UploadImage 上传图片到R2
func (r2 *R2Client) UploadImage(imageData []byte, filename string) (string, error) {
	// 生成唯一的文件名
	uniqueFilename := fmt.Sprintf("bg_removed_%d_%s", time.Now().Unix(), filename)

	// 上传到R2
	_, err := r2.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(r2.config.BucketName),
		Key:         aws.String(uniqueFilename),
		Body:        bytes.NewReader(imageData),
		ContentType: aws.String(getContentType(filename)),
	})
	if err != nil {
		return "", fmt.Errorf("上传到R2失败: %v", err)
	}

	// 构建公开访问URL
	imageURL := r2.buildPublicURL(uniqueFilename)

	return imageURL, nil
}

// buildPublicURL 构建公开访问URL
func (r2 *R2Client) buildPublicURL(filename string) string {
	// 首先检查是否配置了自定义公开域名
	publicDomain := os.Getenv("R2_PUBLIC_DOMAIN")
	if publicDomain != "" {
		return fmt.Sprintf("https://%s/%s", publicDomain, filename)
	}

	// 检查是否配置了R2的公开访问域名
	r2PublicDomain := os.Getenv("R2_DEV_DOMAIN")
	if r2PublicDomain != "" {
		// 如果已经包含 https://，直接使用；否则添加
		if strings.HasPrefix(r2PublicDomain, "https://") {
			return fmt.Sprintf("%s/%s", r2PublicDomain, filename)
		}
		return fmt.Sprintf("https://%s/%s", r2PublicDomain, filename)
	}

	// 如果都没有配置，返回预签名URL（7天有效期）
	presignedURL, err := r2.GeneratePresignedURL(filename, 7*24*time.Hour)
	if err != nil {
		// 如果预签名URL也失败，返回默认格式（可能无法直接访问）
		return fmt.Sprintf("https://%s.%s.r2.cloudflarestorage.com/%s",
			r2.config.BucketName,
			r2.config.AccountID,
			filename)
	}

	return presignedURL
}

// GeneratePresignedURL 生成预签名URL用于访问
func (r2 *R2Client) GeneratePresignedURL(key string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(r2.client)

	presignResult, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(r2.config.BucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})

	if err != nil {
		return "", fmt.Errorf("生成预签名URL失败: %v", err)
	}

	return presignResult.URL, nil
}

// getContentType 根据文件扩展名获取Content-Type
func getContentType(filename string) string {
	if strings.HasSuffix(strings.ToLower(filename), ".png") {
		return "image/png"
	} else if strings.HasSuffix(strings.ToLower(filename), ".jpg") || strings.HasSuffix(strings.ToLower(filename), ".jpeg") {
		return "image/jpeg"
	} else if strings.HasSuffix(strings.ToLower(filename), ".gif") {
		return "image/gif"
	}
	return "image/png" // 默认
}
