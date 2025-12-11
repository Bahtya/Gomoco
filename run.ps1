# Run script for Gomoco

Write-Host "Starting Gomoco..." -ForegroundColor Green

# Check if frontend is built
if (!(Test-Path "web/dist")) {
    Write-Host "Frontend not built. Building now..." -ForegroundColor Yellow
    Set-Location web
    if (!(Test-Path "node_modules")) {
        npm install
    }
    npm run build
    Set-Location ..
}

# Download Go dependencies if needed
if (!(Test-Path "go.sum")) {
    Write-Host "Downloading Go dependencies..." -ForegroundColor Yellow
    go mod download
}

# Run the application
Write-Host "`nStarting Gomoco server on http://localhost:8080" -ForegroundColor Cyan
go run main.go
