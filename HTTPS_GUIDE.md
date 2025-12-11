# HTTPS 使用指南

## 快速开始

### 1. 生成自签名证书（测试用）

**Windows:**
```powershell
# 运行证书生成脚本
.\generate-cert.ps1
```

**Linux/macOS:**
```bash
# 添加执行权限
chmod +x generate-cert.sh

# 运行脚本
./generate-cert.sh

```

这将在 `certs/` 目录下生成：
- `server.crt` - SSL/TLS 证书
- `server.key` - 私钥文件
- `server.pem` - 组合文件（高级脚本）

### 2. 创建 HTTPS Mock API

在 Gomoco Web 界面中：

1. **API 名称**: HTTPS 测试接口
2. **端口**: 9443（或任意端口）
3. **协议**: 选择 **HTTPS**
4. **证书文件路径**: `certs/server.crt`
5. **私钥文件路径**: `certs/server.key`
6. **字符集**: UTF-8
7. **路径**: `/api/test`
8. **方法**: GET
9. **响应内容**: `{"status": "ok", "message": "HTTPS works!"}`

点击"创建 Mock API"。

### 3. 测试 HTTPS Mock API

#### 使用 curl
```bash
# -k 参数跳过证书验证（自签名证书需要）
curl -k https://localhost:9443/api/test
```

#### 使用浏览器
访问 `https://localhost:9443/api/test`

**注意**: 浏览器会显示安全警告（因为是自签名证书），点击"继续访问"或"接受风险"即可。

#### 使用 PowerShell
```powershell
# 跳过证书验证
[System.Net.ServicePointManager]::ServerCertificateValidationCallback = {$true}
Invoke-WebRequest -Uri https://localhost:9443/api/test
```

## 生产环境使用

### 使用 Let's Encrypt 免费证书

```bash
# 1. 安装 certbot
# Ubuntu/Debian
sudo apt-get install certbot

# CentOS/RHEL
sudo yum install certbot

# 2. 获取证书
sudo certbot certonly --standalone -d yourdomain.com

# 3. 证书位置
# 证书: /etc/letsencrypt/live/yourdomain.com/fullchain.pem
# 私钥: /etc/letsencrypt/live/yourdomain.com/privkey.pem

# 4. 在 Gomoco 中配置
# 证书文件: /etc/letsencrypt/live/yourdomain.com/fullchain.pem
# 私钥文件: /etc/letsencrypt/live/yourdomain.com/privkey.pem
```

### 使用商业 CA 证书

如果你购买了商业 SSL 证书：

1. 从 CA 获取证书文件（通常是 `.crt` 或 `.pem` 格式）
2. 获取私钥文件（`.key` 格式）
3. 在 Gomoco 中配置这两个文件的路径

## 证书格式转换

### PFX 转 PEM

如果你有 `.pfx` 或 `.p12` 格式的证书：

```bash
# 提取私钥
openssl pkcs12 -in certificate.pfx -nocerts -out server.key -nodes

# 提取证书
openssl pkcs12 -in certificate.pfx -clcerts -nokeys -out server.crt
```

### DER 转 PEM

```bash
# 证书转换
openssl x509 -inform der -in certificate.cer -out server.crt

# 私钥转换
openssl rsa -inform der -in private.key -out server.key
```



### 参数说明

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-d, --domain` | 主域名 | localhost |
| `-a, --alt-names` | 备用域名（逗号分隔） | 无 |
| `-o, --output` | 输出目录 | certs |
| `-y, --years` | 有效期（年） | 1 |
| `-k, --key-size` | 密钥大小 | 2048 |
| `-h, --help` | 显示帮助 | - |

### 生成的文件

- `server.crt` - X.509 证书
- `server.key` - RSA 私钥
- `server.pem` - 组合文件（证书+私钥）

## 高级配置

### 使用中间证书链

如果你的证书需要中间证书：

```bash
# 合并证书链
cat server.crt intermediate.crt root.crt > fullchain.crt

# 在 Gomoco 中使用 fullchain.crt
```

### 证书权限设置

```bash
# Linux 环境下设置正确的权限
chmod 644 certs/server.crt
chmod 600 certs/server.key

# 确保 Gomoco 进程有读取权限
chown gomoco:gomoco certs/server.key
```

## 常见问题

### Q: 浏览器显示"不安全"警告？
A: 这是因为使用了自签名证书。生产环境请使用 CA 签发的证书。

### Q: 如何信任自签名证书？
A: 
**Windows:**
1. 双击 `server.crt`
2. 点击"安装证书"
3. 选择"本地计算机"
4. 放入"受信任的根证书颁发机构"

**macOS:**
```bash
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain certs/server.crt
```

**Linux:**
```bash
sudo cp certs/server.crt /usr/local/share/ca-certificates/
sudo update-ca-certificates
```

### Q: 证书过期了怎么办？
A: 重新生成证书或续期：
```bash
# 自签名证书
.\generate-cert.ps1

# Let's Encrypt
sudo certbot renew
```

### Q: 可以使用通配符证书吗？
A: 可以！通配符证书（*.example.com）可以用于多个子域名。

### Q: HTTPS 性能如何？
A: Go 的 TLS 实现性能优秀，对于 Mock Server 场景完全够用。

## 安全最佳实践

1. **不要提交私钥到版本控制**
   - 已在 `.gitignore` 中排除 `certs/` 目录
   
2. **使用强密码保护私钥**
   ```bash
   # 生成带密码的私钥
   openssl genrsa -aes256 -out server.key 2048
   ```

3. **定期更新证书**
   - 自签名证书: 每年更新
   - Let's Encrypt: 自动续期（90天）
   - 商业证书: 按购买期限更新

4. **限制私钥文件权限**
   ```bash
   chmod 600 certs/server.key
   ```

5. **使用现代 TLS 版本**
   - Go 默认使用 TLS 1.2 和 1.3
   - 自动禁用不安全的旧版本

## 示例配置

### 开发环境
```yaml
# config/mocks.yaml
mocks:
  - id: "dev-https-api"
    name: "开发环境 HTTPS API"
    port: 9443
    protocol: "https"
    cert_file: "certs/server.crt"
    key_file: "certs/server.key"
    content: '{"env": "development"}'
    charset: "UTF-8"
    path: "/api/status"
    method: "GET"
```

### 生产环境
```yaml
# config/mocks.yaml
mocks:
  - id: "prod-https-api"
    name: "生产环境 HTTPS API"
    port: 443
    protocol: "https"
    cert_file: "/etc/letsencrypt/live/api.example.com/fullchain.pem"
    key_file: "/etc/letsencrypt/live/api.example.com/privkey.pem"
    content: '{"env": "production"}'
    charset: "UTF-8"
    path: "/api/status"
    method: "GET"
```

## 故障排查

### 证书文件找不到
```
HTTPS server error on port 9443: open certs/server.crt: no such file or directory
```
**解决**: 检查证书文件路径是否正确，使用绝对路径或相对于程序运行目录的路径。

### 私钥格式错误
```
HTTPS server error on port 9443: tls: failed to parse private key
```
**解决**: 确保私钥是 PEM 格式，使用 `openssl rsa -in server.key -text` 验证。

### 证书和私钥不匹配
```
HTTPS server error on port 9443: tls: private key does not match public key
```
**解决**: 确保证书和私钥是配对的，重新生成或检查文件。

### 端口权限问题（Linux）
```
HTTPS server error on port 443: bind: permission denied
```
**解决**: 
```bash
# 方法1: 使用 root 运行（不推荐）
sudo ./gomoco-linux-amd64 -port 8080

# 方法2: 授予端口绑定权限
sudo setcap 'cap_net_bind_service=+ep' ./gomoco-linux-amd64

# 方法3: 使用高端口（>1024）
./gomoco-linux-amd64 -port 8080
# 然后用 nginx 反向代理到 443
```

## 参考资源

- [OpenSSL 文档](https://www.openssl.org/docs/)
- [Let's Encrypt](https://letsencrypt.org/)
- [Go TLS 包文档](https://pkg.go.dev/crypto/tls)
- [Mozilla SSL 配置生成器](https://ssl-config.mozilla.org/)
