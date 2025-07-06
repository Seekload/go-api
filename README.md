# Go API 项目

这是一个使用 Go 和 Gin 框架开发的 API 项目，部署在 Vercel 上。

## 功能特性

- **图片生成**：使用第三方API生成吉卜力风格图片
- **图片上传**：支持上传图片到Vercel Blob Storage
- **背景移除**：集成Photoroom API实现智能背景移除功能

## 新增功能：背景移除API

### 概述
新增的背景移除功能使用Photoroom API，能够自动识别并移除图片背景，保留主体对象。据统计，这个功能可以提升电商销售额20-100%。

### API端点
```
POST /api/remove-background
```

### 使用方式

#### 方式1：上传文件
```bash
curl -X POST \
  http://your-domain.com/api/remove-background \
  -H "Content-Type: multipart/form-data" \
  -F "image=@/path/to/your/image.jpg"
```

#### 方式2：使用图片URL
```bash
curl -X POST \
  http://your-domain.com/api/remove-background \
  -H "Content-Type: application/json" \
  -d '{"imageUrl": "https://example.com/image.jpg"}'
```

### 响应格式
```json
{
  "success": true,
  "message": "背景移除成功",
  "imageUrl": "https://your-blob-store.public.blob.vercel-storage.com/bg_removed_1234567890_image.png"
}
```

### 环境变量配置

在使用背景移除功能前，需要配置以下环境变量：

```bash
# Photoroom API密钥 - 从 https://www.photoroom.com/api 获取
PHOTOROOM_API_KEY=your_photoroom_api_key_here

# Vercel Blob存储令牌 - 在Vercel项目的Storage设置中创建
BLOB_READ_WRITE_TOKEN=your_blob_token_here

# Vercel Blob存储URL - 你的Blob Store的URL
VERCEL_BLOB_STORE_URL=https://your-store-id.public.blob.vercel-storage.com
```

#### 获取Photoroom API密钥
1. 访问 [Photoroom API](https://www.photoroom.com/api)
2. 注册账户并登录
3. 在API控制台中创建新的API密钥
4. 复制密钥并设置为环境变量

#### 配置Vercel Blob Storage
1. 在Vercel项目中进入Storage标签
2. 创建新的Blob Store或使用现有的
3. 复制Read-Write Token
4. 记录Blob Store的URL（格式：https://store-id.public.blob.vercel-storage.com）

### 技术实现细节

#### 工作流程
1. **接收请求**：支持文件上传或图片URL
2. **调用Photoroom API**：使用multipart/form-data格式发送图片数据
3. **处理响应**：接收处理后的图片数据
4. **存储结果**：将处理后的图片上传到Vercel Blob Storage
5. **返回URL**：返回可访问的图片URL

#### 支持的图片格式
- JPEG (.jpg, .jpeg)
- PNG (.png)  
- GIF (.gif)

#### 性能特点
- **智能处理**：AI驱动的背景移除，效果优于传统方法
- **全球CDN**：通过Vercel Blob Storage提供全球访问
- **高可用性**：99.99%的可用性保证

### 错误处理

常见错误及解决方案：

```json
// 未配置API密钥
{
  "success": false,
  "message": "背景移除失败: 未配置Photoroom API密钥"
}

// 无效的图片格式
{
  "success": false,
  "message": "只支持 JPG、PNG、GIF 格式的图片"
}

// 网络错误
{
  "success": false,
  "message": "背景移除失败: API请求失败，状态码: 400"
}
```

## 部署配置

### Vercel部署
确保在Vercel项目设置中添加所有必需的环境变量：

1. 进入Vercel项目设置
2. 选择Environment Variables
3. 添加以下变量：
   - `PHOTOROOM_API_KEY`
   - `BLOB_READ_WRITE_TOKEN`
   - `VERCEL_BLOB_STORE_URL`

### 本地开发
```bash
# 拉取环境变量
vercel env pull

# 启动开发服务器
vercel dev
```

## API文档

### 现有端点
- `GET /api/hello` - 测试端点
- `GET /api/ping` - 健康检查
- `POST /api/generate-image` - 生成图片
- `GET /api/getTaskInfo` - 获取任务信息
- `POST /api/uploadImg` - 上传图片

### 新增端点
- `POST /api/remove-background` - 移除图片背景

## 开发规范

遵循项目的.cursorrules文件中定义的Go代码规范和注释标准。

## 许可证

请查看LICENSE文件了解详细信息。 