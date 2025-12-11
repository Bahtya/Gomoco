# Cross-compile script for Linux (from Windows)

Write-Host "Building Gomoco for Linux..." -ForegroundColor Green

# Build frontend first
Write-Host "`nBuilding frontend..." -ForegroundColor Yellow
Set-Location web
if (!(Test-Path "node_modules")) {
    Write-Host "Installing frontend dependencies..." -ForegroundColor Yellow
    npm install
}
npm run build
Set-Location ..

# Set environment variables for Linux cross-compilation
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

Write-Host "`nCross-compiling for Linux (amd64) with static linking..." -ForegroundColor Yellow
go mod download
go build -ldflags="-s -w -extldflags '-static'" -tags netgo -a -o gomoco-linux-amd64 main.go

# Reset environment variables
Remove-Item Env:\GOOS
Remove-Item Env:\GOARCH
Remove-Item Env:\CGO_ENABLED

Write-Host "`nBuild completed successfully!" -ForegroundColor Green
Write-Host "Output: gomoco-linux-amd64" -ForegroundColor Cyan
Write-Host "`nTo run on Linux:" -ForegroundColor Yellow
Write-Host "  chmod +x gomoco-linux-amd64" -ForegroundColor White
Write-Host "  ./gomoco-linux-amd64" -ForegroundColor White
