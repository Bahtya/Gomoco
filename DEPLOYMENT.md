# Gomoco 部署指南

## 快速部署

Gomoco 使用单一可执行文件部署，无需额外依赖。

## 构建

### Windows 环境

#### 本地构建
```powershell
.\build.ps1
```
生成 `gomoco.exe`

#### 交叉编译 Linux
```powershell
.\build-linux.ps1
```
生成 `gomoco-linux-amd64`

#### 构建所有平台
```powershell
.\build-all.ps1
```
生成所有平台的可执行文件

## 部署到 Linux 服务器

### 方式 1: 直接运行

```bash
# 1. 上传文件
scp gomoco-linux-amd64 user@server:/opt/gomoco/

# 2. 添加执行权限
ssh user@server
cd /opt/gomoco
chmod +x gomoco-linux-amd64

# 3. 运行（默认端口 8080）
./gomoco-linux-amd64

# 3a. 自定义端口运行
./gomoco-linux-amd64 -port 9000

# 3b. 查看版本
./gomoco-linux-amd64 -version
```

### 方式 2: 使用 systemd 服务

创建服务文件 `/etc/systemd/system/gomoco.service`:

```ini
[Unit]
Description=Gomoco Mock Server
After=network.target

[Service]
Type=simple
User=gomoco
WorkingDirectory=/opt/gomoco
ExecStart=/opt/gomoco/gomoco-linux-amd64 -port 8080
Restart=on-failure
RestartSec=5s

# 环境变量（可选）
# Environment="GOMOCO_PORT=8080"

[Install]
WantedBy=multi-user.target
```

启动服务:
```bash
sudo systemctl daemon-reload
sudo systemctl enable gomoco
sudo systemctl start gomoco
sudo systemctl status gomoco
```

### 方式 3: 使用 Docker

#### 选项 A: 最小镜像（推荐）
由于使用了静态编译，可以使用 `scratch` 基础镜像：

```dockerfile
FROM scratch
COPY gomoco-linux-amd64 /gomoco
EXPOSE 8080
CMD ["/gomoco"]
```

镜像大小仅为可执行文件大小！

#### 选项 B: Alpine 镜像
如果需要 shell 调试：

```dockerfile
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY gomoco-linux-amd64 /app/gomoco
RUN chmod +x /app/gomoco
EXPOSE 8080
CMD ["/app/gomoco"]
```

构建和运行:
```bash
docker build -t gomoco:latest .
docker run -d -p 8080:8080 -v $(pwd)/config:/app/config --name gomoco gomoco:latest
```

### 方式 4: 使用 Docker Compose

创建 `docker-compose.yml`:
```yaml
version: '3.8'
services:
  gomoco:
    image: gomoco:latest
    container_name: gomoco
    ports:
      - "8080:8080"
    volumes:
      - ./config:/app/config
    restart: unless-stopped
```

运行:
```bash
docker-compose up -d
```

## 配置持久化

配置文件位置: `config/mocks.yaml`

### 备份配置
```bash
cp config/mocks.yaml config/mocks.yaml.backup
```

### 恢复配置
```bash
cp config/mocks.yaml.backup config/mocks.yaml
```

### 迁移到新服务器
只需复制 `config/mocks.yaml` 文件到新服务器的 `config` 目录。

## 端口说明

- **8080**: API 服务器和 Web 界面（默认，可通过 `-port` 参数修改）
- **自定义端口**: 用户创建的 Mock API 端口

**修改 API 服务器端口：**
```bash
# 使用 9000 端口
./gomoco-linux-amd64 -port 9000

# 访问地址变为
# http://localhost:9000
```

确保防火墙允许这些端口的访问。

## 日志管理

### 查看日志（systemd）
```bash
sudo journalctl -u gomoco -f
```

### 查看日志（直接运行）
```bash
nohup ./gomoco-linux-amd64 > gomoco.log 2>&1 &
tail -f gomoco.log
```

## 性能优化

1. **编译优化**: 
   - 使用 `-ldflags="-s -w"` 减小文件大小
   - Linux 版本使用静态链接，无运行时依赖
   - 使用 `-tags netgo` 纯 Go 网络库
2. **内存使用**: 每个 Mock API 占用极少内存
3. **并发处理**: Go 原生支持高并发
4. **启动速度**: 静态编译版本启动更快

## 安全建议

1. 使用反向代理（Nginx/Caddy）添加 HTTPS
2. 限制管理界面访问（通过防火墙或反向代理）
3. 定期备份配置文件
4. 使用非 root 用户运行

## 故障排查

### 端口被占用
```bash
# 查看端口占用
netstat -tlnp | grep :8080

# 或使用 ss
ss -tlnp | grep :8080
```

### 权限问题
```bash
# 确保可执行权限
chmod +x gomoco-linux-amd64

# 确保配置目录权限
chmod 755 config
chmod 644 config/mocks.yaml
```

### 服务无法启动
```bash
# 检查日志
journalctl -u gomoco -n 50

# 手动运行查看错误
./gomoco-linux-amd64
```

## 更新部署

1. 构建新版本
2. 停止旧服务
3. 替换可执行文件
4. 启动新服务

```bash
# systemd 方式
sudo systemctl stop gomoco
sudo cp gomoco-linux-amd64 /opt/gomoco/
sudo systemctl start gomoco

# Docker 方式
docker-compose down
docker build -t gomoco:latest .
docker-compose up -d
```

配置文件会自动保留和加载。
