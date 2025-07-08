package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go-api/api/controllers"
)

func main() {
	fmt.Println("🚀 R2上传测试Demo")
	fmt.Println("==================")

	// 检查本地图片文件是否存在
	imagePath := "images/image1.jpg"
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Fatalf("❌ 图片文件不存在: %s", imagePath)
	}

	fmt.Printf("📁 准备上传文件: %s\n", imagePath)

	// 读取图片文件
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		log.Fatalf("❌ 读取图片文件失败: %v", err)
	}

	fmt.Printf("📊 文件大小: %.2f KB\n", float64(len(imageData))/1024)

	// 创建R2客户端
	fmt.Println("\n🔗 连接到Cloudflare R2...")
	r2Client, err := controllers.NewR2Client()
	if err != nil {
		log.Fatalf("❌ 创建R2客户端失败: %v", err)
	}

	fmt.Println("✅ R2客户端创建成功")

	// 上传图片到R2
	fmt.Println("\n📤 开始上传图片到R2...")
	filename := filepath.Base(imagePath)
	imageURL, err := r2Client.UploadImage(imageData, filename)
	if err != nil {
		log.Fatalf("❌ 上传失败: %v", err)
	}

	// 显示结果
	fmt.Println("\n🎉 上传成功！")
	fmt.Printf("🔗 图片URL: %s\n", imageURL)
	fmt.Printf("📝 文件名: %s\n", filename)
	fmt.Printf("📏 文件大小: %d bytes\n", len(imageData))

	fmt.Println("\n✨ 测试完成！")
	fmt.Println("💡 提示：如果你的R2存储桶设置为公开访问，可以直接在浏览器中打开上述URL查看图片。")
}
