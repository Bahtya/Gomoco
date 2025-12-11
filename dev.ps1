# Development script for Gomoco

Write-Host "Starting Gomoco in development mode..." -ForegroundColor Green

# Start backend in background
Write-Host "`nStarting backend server..." -ForegroundColor Yellow
$backend = Start-Process powershell -ArgumentList "-NoExit", "-Command", "go run main.go" -PassThru

# Wait a bit for backend to start
Start-Sleep -Seconds 2

# Start frontend dev server
Write-Host "Starting frontend dev server..." -ForegroundColor Yellow
Set-Location web
if (!(Test-Path "node_modules")) {
    Write-Host "Installing frontend dependencies..." -ForegroundColor Yellow
    npm install
}

Write-Host "`nFrontend: http://localhost:3000" -ForegroundColor Cyan
Write-Host "Backend: http://localhost:8080" -ForegroundColor Cyan
Write-Host "`nPress Ctrl+C to stop both servers" -ForegroundColor Yellow

npm run dev

# Cleanup
Set-Location ..
Stop-Process -Id $backend.Id -Force
