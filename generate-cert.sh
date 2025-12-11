#!/bin/bash
# 生成自签名 SSL/TLS 证书用于 HTTPS 测试
# Generate self-signed SSL/TLS certificate for HTTPS testing

set -e

echo "Generating self-signed certificate for HTTPS testing..."

# 创建证书目录
CERT_DIR="certs"
if [ ! -d "$CERT_DIR" ]; then
    mkdir -p "$CERT_DIR"
    echo "✓ Created directory: $CERT_DIR"
fi

# 检查 OpenSSL 是否安装
if ! command -v openssl &> /dev/null; then
    echo "Error: OpenSSL is not installed."
    echo "Please install OpenSSL first:"
    echo "  Ubuntu/Debian: sudo apt-get install openssl"
    echo "  CentOS/RHEL:   sudo yum install openssl"
    echo "  Alpine:        apk add openssl"
    exit 1
fi

echo ""
echo "Using OpenSSL to generate certificate..."

# 证书配置
COUNTRY="CN"
STATE="Beijing"
CITY="Beijing"
ORGANIZATION="Gomoco"
ORGANIZATIONAL_UNIT="Dev"
COMMON_NAME="localhost"
DAYS=365

# 生成私钥
echo "→ Generating private key..."
openssl genrsa -out "$CERT_DIR/server.key" 2048

# 生成证书签名请求（CSR）
echo "→ Generating certificate signing request..."
openssl req -new \
    -key "$CERT_DIR/server.key" \
    -out "$CERT_DIR/server.csr" \
    -subj "/C=$COUNTRY/ST=$STATE/L=$CITY/O=$ORGANIZATION/OU=$ORGANIZATIONAL_UNIT/CN=$COMMON_NAME"

# 生成自签名证书（有效期 365 天）
echo "→ Generating self-signed certificate..."
openssl x509 -req \
    -days $DAYS \
    -in "$CERT_DIR/server.csr" \
    -signkey "$CERT_DIR/server.key" \
    -out "$CERT_DIR/server.crt"

# 删除 CSR 文件
rm -f "$CERT_DIR/server.csr"

# 设置正确的权限
chmod 644 "$CERT_DIR/server.crt"
chmod 600 "$CERT_DIR/server.key"

echo ""
echo "========================================"
echo "✓ Certificate generated successfully!"
echo "========================================"
echo "Certificate: $CERT_DIR/server.crt"
echo "Private Key: $CERT_DIR/server.key"
echo ""
echo "Certificate Details:"
openssl x509 -in "$CERT_DIR/server.crt" -noout -subject -dates

echo ""
echo "========================================"
echo "Usage in Gomoco:"
echo "========================================"
echo "1. Create a Mock API with protocol: HTTPS"
echo "2. Certificate file: $CERT_DIR/server.crt"
echo "3. Key file: $CERT_DIR/server.key"
echo ""
echo "Test with curl:"
echo "  curl -k https://localhost:9443/api/test"
echo ""
echo "Note: Self-signed certificates will show"
echo "security warnings in browsers."
echo "For production, use certificates from a"
echo "trusted CA (e.g., Let's Encrypt)."
echo "========================================"
