# Gomoco 使用示例

## 命令行参数

### 基本使用

```bash
# 默认端口 8080 启动
./gomoco

# 输出:
# Starting Gomoco v1.1.0 on http://localhost:8080
```

### 自定义端口

```bash
# 使用 9000 端口
./gomoco -port 9000

# 输出:
# Starting Gomoco v1.1.0 on http://localhost:9000
```

### 查看版本

```bash
./gomoco -version

# 输出:
# Gomoco v1.1.0
# A lightweight mock server written in Go
```

### 查看帮助

```bash
./gomoco -h

# 输出:
# Usage of ./gomoco:
#   -port int
#         API server port (default 8080)
#   -version
#         Show version information
```

## 不同场景使用

### 场景 1: 开发环境

```bash
# 开发时使用默认端口
./gomoco

# 访问管理界面
# http://localhost:8080
```

### 场景 2: 多实例运行

```bash
# 实例 1 - 测试环境
./gomoco -port 8080 &

# 实例 2 - 预发布环境
./gomoco -port 8081 &

# 实例 3 - 演示环境
./gomoco -port 8082 &
```

### 场景 3: 生产环境

```bash
# 后台运行，指定端口
nohup ./gomoco -port 8080 > gomoco.log 2>&1 &

# 查看进程
ps aux | grep gomoco

# 查看日志
tail -f gomoco.log
```

### 场景 4: 反向代理后

```bash
# Gomoco 运行在内网端口
./gomoco -port 8888

# Nginx 配置
# server {
#     listen 80;
#     server_name mock.example.com;
#     location / {
#         proxy_pass http://localhost:8888;
#     }
# }
```

## Docker 使用

### 默认端口

```bash
# 使用默认端口 8080
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/config:/config \
  --name gomoco \
  gomoco:latest
```

### 自定义端口

```bash
# 使用 9000 端口
docker run -d \
  -p 9000:9000 \
  -v $(pwd)/config:/config \
  --name gomoco \
  gomoco:latest \
  /gomoco -port 9000
```

### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'
services:
  gomoco:
    image: gomoco:latest
    command: ["/gomoco", "-port", "9000"]
    ports:
      - "9000:9000"
    volumes:
      - ./config:/config
```

```bash
docker-compose up -d
```

## Systemd 服务

### 默认端口

```ini
# /etc/systemd/system/gomoco.service
[Unit]
Description=Gomoco Mock Server
After=network.target

[Service]
Type=simple
User=gomoco
WorkingDirectory=/opt/gomoco
ExecStart=/opt/gomoco/gomoco-linux-amd64
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

### 自定义端口

```ini
# /etc/systemd/system/gomoco.service
[Unit]
Description=Gomoco Mock Server
After=network.target

[Service]
Type=simple
User=gomoco
WorkingDirectory=/opt/gomoco
ExecStart=/opt/gomoco/gomoco-linux-amd64 -port 9000
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable gomoco
sudo systemctl start gomoco
```

## 端口冲突处理

### 检查端口占用

```bash
# Linux
netstat -tlnp | grep :8080
# 或
ss -tlnp | grep :8080

# Windows
netstat -ano | findstr :8080
```

### 使用其他端口

```bash
# 如果 8080 被占用，使用其他端口
./gomoco -port 8888
```

## 防火墙配置

### Linux (iptables)

```bash
# 允许 8080 端口
sudo iptables -A INPUT -p tcp --dport 8080 -j ACCEPT
```

### Linux (firewalld)

```bash
# 允许 8080 端口
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --reload
```

### Linux (ufw)

```bash
# 允许 8080 端口
sudo ufw allow 8080/tcp
```

## 性能测试

### 使用 Apache Bench

```bash
# 创建一个 Mock API (端口 9090)
# 然后测试性能

ab -n 10000 -c 100 http://localhost:9090/api/test
```

### 使用 wrk

```bash
wrk -t4 -c100 -d30s http://localhost:9090/api/test
```

## 监控和日志

### 查看实时日志

```bash
# 直接运行
tail -f gomoco.log

# Systemd
journalctl -u gomoco -f

# Docker
docker logs -f gomoco
```

### 日志轮转

```bash
# /etc/logrotate.d/gomoco
/opt/gomoco/gomoco.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0644 gomoco gomoco
}
```

## 常见问题

### Q: 如何更改默认端口？
A: 使用 `-port` 参数：
```bash
./gomoco -port 9000
```

### Q: 可以同时运行多个实例吗？
A: 可以，只要使用不同的端口：
```bash
./gomoco -port 8080 &
./gomoco -port 8081 &
```

### Q: 如何在后台运行？
A: 使用 nohup 或 systemd：
```bash
nohup ./gomoco -port 8080 > gomoco.log 2>&1 &
```

### Q: 如何停止服务？
A: 
```bash
# 查找进程
ps aux | grep gomoco

# 停止进程
kill <PID>

# 或使用 systemd
sudo systemctl stop gomoco
```

### Q: 端口被占用怎么办？
A: 使用其他端口或停止占用端口的程序：
```bash
# 查看占用
netstat -tlnp | grep :8080

# 使用其他端口
./gomoco -port 8888
```

## 最佳实践

1. **生产环境使用 systemd**
   - 自动重启
   - 日志管理
   - 开机自启

2. **使用反向代理**
   - Nginx/Caddy 添加 HTTPS
   - 负载均衡
   - 访问控制

3. **定期备份配置**
   ```bash
   cp config/mocks.yaml config/mocks.yaml.backup
   ```

4. **监控端口使用**
   - 避免端口冲突
   - 记录 Mock API 端口

5. **日志管理**
   - 使用 logrotate
   - 定期清理旧日志
