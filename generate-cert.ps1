# 生成自签名 SSL/TLS 证书用于 HTTPS 测试
# Generate self-signed SSL/TLS certificate for HTTPS testing

Write-Host "Generating self-signed certificate for HTTPS testing..." -ForegroundColor Yellow

# 创建证书目录
$certDir = "certs"
if (!(Test-Path $certDir)) {
    New-Item -ItemType Directory -Path $certDir | Out-Null
    Write-Host "Created directory: $certDir" -ForegroundColor Green
}

# 使用 OpenSSL 生成证书（如果已安装）
if (Get-Command openssl -ErrorAction SilentlyContinue) {
    Write-Host "`nUsing OpenSSL to generate certificate..." -ForegroundColor Cyan
    
    # 生成私钥
    openssl genrsa -out "$certDir/server.key" 2048
    
    # 生成证书签名请求（CSR）
    openssl req -new -key "$certDir/server.key" -out "$certDir/server.csr" -subj "/C=CN/ST=Beijing/L=Beijing/O=Gomoco/OU=Dev/CN=localhost"
    
    # 生成自签名证书（有效期 365 天）
    openssl x509 -req -days 365 -in "$certDir/server.csr" -signkey "$certDir/server.key" -out "$certDir/server.crt"
    
    # 删除 CSR 文件
    Remove-Item "$certDir/server.csr" -ErrorAction SilentlyContinue
    
    Write-Host "`nCertificate generated successfully!" -ForegroundColor Green
    Write-Host "Certificate: $certDir/server.crt" -ForegroundColor Cyan
    Write-Host "Private Key: $certDir/server.key" -ForegroundColor Cyan
    
} else {
    # 使用 PowerShell 生成证书（Windows 内置方法）
    Write-Host "`nUsing PowerShell to generate certificate..." -ForegroundColor Cyan
    
    $cert = New-SelfSignedCertificate `
        -DnsName "localhost", "127.0.0.1" `
        -CertStoreLocation "Cert:\CurrentUser\My" `
        -NotAfter (Get-Date).AddYears(1) `
        -KeyAlgorithm RSA `
        -KeyLength 2048 `
        -KeyUsage DigitalSignature, KeyEncipherment `
        -TextExtension @("2.5.29.37={text}1.3.6.1.5.5.7.3.1")
    
    # 导出证书
    $certPath = "$certDir/server.crt"
    $keyPath = "$certDir/server.key"
    $pfxPath = "$certDir/server.pfx"
    $password = ConvertTo-SecureString -String "gomoco" -Force -AsPlainText
    
    # 导出为 PFX
    Export-PfxCertificate -Cert $cert -FilePath $pfxPath -Password $password | Out-Null
    
    # 导出证书（CRT）
    Export-Certificate -Cert $cert -FilePath $certPath -Type CERT | Out-Null
    
    Write-Host "`nCertificate generated successfully!" -ForegroundColor Green
    Write-Host "Certificate: $certPath" -ForegroundColor Cyan
    Write-Host "PFX File: $pfxPath (password: gomoco)" -ForegroundColor Cyan
    Write-Host "`nNote: To use with Gomoco, you need to convert PFX to PEM format:" -ForegroundColor Yellow
    Write-Host "  openssl pkcs12 -in $pfxPath -nocerts -out $keyPath -nodes -passin pass:gomoco" -ForegroundColor White
    Write-Host "  openssl pkcs12 -in $pfxPath -clcerts -nokeys -out $certPath -passin pass:gomoco" -ForegroundColor White
}

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "Usage in Gomoco:" -ForegroundColor Yellow
Write-Host "1. Create a Mock API with protocol: HTTPS" -ForegroundColor White
Write-Host "2. Certificate file: certs/server.crt" -ForegroundColor White
Write-Host "3. Key file: certs/server.key" -ForegroundColor White
Write-Host "`nNote: Self-signed certificates will show security warnings in browsers." -ForegroundColor Yellow
Write-Host "For production, use certificates from a trusted CA." -ForegroundColor Yellow
Write-Host "========================================`n" -ForegroundColor Cyan
