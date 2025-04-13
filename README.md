# Go Server on Vercel

这是一个简单的 Go 语言服务器，部署在 Vercel 上。

## 部署步骤

1. 确保你已经安装了 [Vercel CLI](https://vercel.com/docs/cli)
2. 登录 Vercel：
   ```bash
   vercel login
   ```
3. 部署项目：
   ```bash
   vercel
   ```

## 本地开发

1. 安装依赖：
   ```bash
   go mod tidy
   ```
2. 运行服务器：
   ```bash
   go run main.go
   ```
3. 访问 http://localhost:8080 查看效果 