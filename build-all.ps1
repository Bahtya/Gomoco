# Build script for all platforms

Write-Host "Building Gomoco for all platforms..." -ForegroundColor Green

# Build frontend first
Write-Host "`nBuilding frontend..." -ForegroundColor Yellow
Set-Location web
if (!(Test-Path "node_modules")) {
    Write-Host "Installing frontend dependencies..." -ForegroundColor Yellow
    npm install
}
npm run build
Set-Location ..

# Download Go dependencies
Write-Host "`nDownloading Go dependencies..." -ForegroundColor Yellow
go mod download

# Build for Windows
Write-Host "`nBuilding for Windows (amd64)..." -ForegroundColor Yellow
$env:GOOS = "windows"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -ldflags="-s -w" -tags netgo -a -o gomoco-windows-amd64.exe main.go

# Build for Linux (amd64) - Static linking
Write-Host "Building for Linux (amd64) with static linking..." -ForegroundColor Yellow
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -ldflags="-s -w -extldflags '-static'" -tags netgo -a -o gomoco-linux-amd64 main.go

# Build for Linux ARM64 - Static linking
Write-Host "Building for Linux (arm64) with static linking..." -ForegroundColor Yellow
$env:GOARCH = "arm64"
go build -ldflags="-s -w -extldflags '-static'" -tags netgo -a -o gomoco-linux-arm64 main.go

# Build for macOS (amd64)
Write-Host "Building for macOS (amd64)..." -ForegroundColor Yellow
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -ldflags="-s -w" -tags netgo -a -o gomoco-darwin-amd64 main.go

# Build for macOS ARM64 (Apple Silicon)
Write-Host "Building for macOS (arm64)..." -ForegroundColor Yellow
$env:GOARCH = "arm64"
go build -ldflags="-s -w" -tags netgo -a -o gomoco-darwin-arm64 main.go

# Reset environment variables
Remove-Item Env:\GOOS
Remove-Item Env:\GOARCH
Remove-Item Env:\CGO_ENABLED

Write-Host "`nBuild completed successfully!" -ForegroundColor Green
Write-Host "`nGenerated binaries:" -ForegroundColor Cyan
Write-Host "  - gomoco-windows-amd64.exe (Windows 64-bit)" -ForegroundColor White
Write-Host "  - gomoco-linux-amd64       (Linux 64-bit)" -ForegroundColor White
Write-Host "  - gomoco-linux-arm64       (Linux ARM64)" -ForegroundColor White
Write-Host "  - gomoco-darwin-amd64      (macOS Intel)" -ForegroundColor White
Write-Host "  - gomoco-darwin-arm64      (macOS Apple Silicon)" -ForegroundColor White
