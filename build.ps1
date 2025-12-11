# Build script for Gomoco

Write-Host "Building Gomoco..." -ForegroundColor Green

# Build frontend
Write-Host "`nBuilding frontend..." -ForegroundColor Yellow
Set-Location web
if (!(Test-Path "node_modules")) {
    Write-Host "Installing frontend dependencies..." -ForegroundColor Yellow
    npm install
}
npm run build
Set-Location ..

# Build backend
Write-Host "`nBuilding backend..." -ForegroundColor Yellow
$env:CGO_ENABLED = "0"
go mod download
go build -ldflags="-s -w" -tags netgo -a -o gomoco.exe main.go
Remove-Item Env:\CGO_ENABLED -ErrorAction SilentlyContinue

Write-Host "`nBuild completed successfully!" -ForegroundColor Green
Write-Host "Run './gomoco.exe' to start the server" -ForegroundColor Cyan
Write-Host "`nNote: The executable is self-contained with embedded web assets" -ForegroundColor Yellow
