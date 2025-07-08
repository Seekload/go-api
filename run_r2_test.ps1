# R2上传测试脚本

Write-Host "🚀 R2上传测试准备" -ForegroundColor Green
Write-Host "==================" -ForegroundColor Green

# 检查环境变量是否设置
$envVars = @("R2_ACCOUNT_ID", "R2_ACCESS_KEY_ID", "R2_SECRET_ACCESS_KEY", "R2_BUCKET_NAME")
$missingVars = @()

foreach ($var in $envVars) {
    if (-not [Environment]::GetEnvironmentVariable($var)) {
        $missingVars += $var
    }
}

if ($missingVars.Count -gt 0) {
    Write-Host "❌ 缺少以下环境变量:" -ForegroundColor Red
    foreach ($var in $missingVars) {
        Write-Host "   - $var" -ForegroundColor Yellow
    }
    Write-Host ""
    Write-Host "💡 解决方案:" -ForegroundColor Cyan
    Write-Host "1. 复制 r2_config_example.txt 为 .env.local" -ForegroundColor White
    Write-Host "2. 在 .env.local 中填入你的实际配置" -ForegroundColor White
    Write-Host "3. 运行: Get-Content .env.local | ForEach-Object { if(`$_ -match '^(.+)=(.+)$') { [Environment]::SetEnvironmentVariable(`$matches[1], `$matches[2]) } }" -ForegroundColor White
    Write-Host "4. 重新运行此脚本" -ForegroundColor White
    exit 1
}

Write-Host "✅ 环境变量检查通过" -ForegroundColor Green

# 检查图片文件是否存在
if (-not (Test-Path "images/image1.jpg")) {
    Write-Host "❌ 测试图片不存在: images/image1.jpg" -ForegroundColor Red
    exit 1
}

Write-Host "✅ 测试图片存在" -ForegroundColor Green
Write-Host ""

# 运行测试
Write-Host "🔥 开始运行R2上传测试..." -ForegroundColor Magenta
go run test_r2_upload.go 