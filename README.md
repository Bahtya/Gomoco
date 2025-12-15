# Gomoco - Go Mock Server

Gomoco 是一个轻量级的 Mock Server 工具，灵感来源于 [Moco](https://github.com/dreamhead/moco)。它提供了简单易用的 Web 界面来管理和配置 Mock API。

## 功能特性

- ✅ 支持 HTTP、HTTPS、TCP 和 FTP 协议
- ✅ HTTPS 支持自定义 SSL/TLS 证书
- ✅ FTP 支持主动/被动模式，Web 端文件管理
- ✅ 文件上传限制 100MB
- ✅ 支持 UTF-8 和 GBK 字符集编码
- ✅ 固定报文内容响应
- ✅ 动态端口配置
- ✅ HTTP 路径和方法配置
- ✅ API 名称管理
- ✅ **配置持久化** (YAML 文件存储)
- ✅ 自动恢复已保存的 Mock API
- ✅ **单一可执行文件部署** (嵌入式前端资源)
- ✅ **跨平台支持** (Windows/Linux/macOS)
- ✅ 现代化 Vue 前端界面
- ✅ RESTful API 管理接口

## 技术栈

### 后端
- Go 1.21+
- Gin Web Framework
- golang.org/x/text (字符集转换)

### 前端
- Vue 3
- Vite
- Axios

## 快速开始

### 前置要求

- Go 1.21 或更高版本
- Node.js 16+ 和 npm (用于前端开发)

### 安装依赖

#### 后端依赖
```bash
go mod download
```

#### 前端依赖
```bash
cd web
npm install
```

### 开发模式

#### 1. 构建前端
```bash
cd web
npm run build
cd ..
```

#### 2. 启动后端服务
```bash
# 默认端口 8080
go run main.go

# 自定义端口
go run main.go -port 9000

# 查看版本
go run main.go -version

# 查看帮助
go run main.go -h
```

服务默认在 `http://localhost:8080` 启动。

#### 前端开发模式（可选）
如果需要修改前端代码，可以使用开发模式：
```bash
cd web
npm run dev
```

前端开发服务器将在 `http://localhost:3000` 启动，并自动代理 API 请求到后端。

### 构建生产版本

#### Windows 本地构建
```powershell
.\build.ps1
```
生成 `gomoco.exe`，这是一个**自包含的单一可执行文件**，包含了所有前端资源。

#### 交叉编译 Linux 版本（从 Windows）
```powershell
.\build-linux.ps1
```
生成 `gomoco-linux-amd64`，可直接在 Linux 系统上运行。

#### 构建所有平台版本
```powershell
.\build-all.ps1
```
生成以下可执行文件：
- `gomoco-windows-amd64.exe` - Windows 64位
- `gomoco-linux-amd64` - Linux 64位
- `gomoco-linux-arm64` - Linux ARM64
- `gomoco-darwin-amd64` - macOS Intel
- `gomoco-darwin-arm64` - macOS Apple Silicon

**编译选项说明：**
- `-ldflags="-s -w"` - 去除调试信息，减小文件大小
- `-ldflags="-extldflags '-static'"` - Linux 静态链接，无需依赖系统库
- `-tags netgo` - 使用纯 Go 网络库，避免 CGO 依赖
- `-a` - 强制重新编译所有包
- `CGO_ENABLED=0` - 禁用 CGO，确保完全静态编译

**静态编译优势：**
- ✅ 无需安装任何运行时依赖
- ✅ 可在任何 Linux 发行版上运行（包括 Alpine、BusyBox）
- ✅ 适合容器化部署（FROM scratch）
- ✅ 避免 glibc 版本兼容问题

## 命令行参数

```bash
gomoco [选项]

选项:
  -port int
        API 服务器端口 (默认: 8080)
  -version
        显示版本信息
  -h, -help
        显示帮助信息
```

**使用示例：**
```bash
# 默认端口 8080
./gomoco

# 自定义端口
./gomoco -port 9000

# 查看版本
./gomoco -version
```

## 使用说明

### 创建 Mock API

1. 打开浏览器访问 `http://localhost:8080` (或自定义端口)
2. 填写表单：
   - **API 名称**: 给 Mock API 起一个描述性的名称
   - **端口**: Mock 服务监听的端口 (1-65535)
   - **协议**: HTTP、HTTPS、TCP 或 FTP
   - **字符集**: UTF-8 或 GBK
   - **响应内容**: 固定返回的报文内容（FTP 不需要）
   - **路径** (HTTP/HTTPS): HTTP 请求路径，默认为 `/`
   - **方法** (HTTP/HTTPS): HTTP 方法，留空表示任意方法
   - **证书文件** (HTTPS): SSL/TLS 证书文件路径
   - **私钥文件** (HTTPS): SSL/TLS 私钥文件路径
   - **FTP 模式** (FTP): 主动模式或被动模式
   - **根目录** (FTP): FTP 文件根目录，留空自动生成
   - **用户名** (FTP): FTP 登录用户名，默认 admin
   - **密码** (FTP): FTP 登录密码，默认 admin
   - **被动端口范围** (FTP): 被动模式端口范围，如 50000-50100
3. 点击"创建 Mock API"

**注意**: 所有配置会自动保存到 `config/mocks.yaml` 文件中，重启后自动恢复。

### 编辑 Mock API

1. 在列表中找到要编辑的 Mock API
2. 点击"编辑"按钮
3. 修改内容后点击"更新 Mock API"

### 删除 Mock API

1. 在列表中找到要删除的 Mock API
2. 点击"删除"按钮
3. 确认删除

### 测试 Mock API

#### HTTP 示例
```bash
# 假设创建了一个 HTTP Mock API，端口 9090，路径 /test
curl http://localhost:9090/test
```

#### HTTPS 示例
```bash
# 1. 生成自签名证书
# Windows:
.\generate-cert.ps1

# Linux/macOS:
chmod +x generate-cert.sh
./generate-cert.sh

# 2. 创建 HTTPS Mock API（使用生成的证书）
# - 协议: HTTPS
# - 端口: 9443
# - 证书文件: certs/server.crt
# - 私钥文件: certs/server.key

# 3. 测试（-k 跳过证书验证）
curl -k https://localhost:9443/test

# 或使用浏览器访问（会显示安全警告，点击继续即可）
# https://localhost:9443/test
```

#### TCP 示例
```bash
# 假设创建了一个 TCP Mock API，端口 9091
echo "test" | nc localhost 9091
```

#### FTP 示例
```bash
# 1. 在 Web 界面创建 FTP Mock API
# - 协议: FTP
# - 端口: 21
# - FTP 模式: 被动模式
# - 用户名: admin
# - 密码: admin

# 2. 使用 FTP 客户端连接
ftp localhost 21
# 输入用户名: admin
# 输入密码: admin

# 3. 或使用命令行工具
curl -u admin:admin ftp://localhost:21/

# 4. 上传文件
curl -u admin:admin -T myfile.txt ftp://localhost:21/

# 5. 下载文件
curl -u admin:admin -O ftp://localhost:21/myfile.txt

# 6. 使用 Web API 管理文件
# 列出文件
curl http://localhost:8080/api/mocks/{mock_id}/files

# 上传文件（限制 100MB）
curl -F "file=@myfile.txt" -F "path=/" http://localhost:8080/api/mocks/{mock_id}/files

# 下载文件
curl http://localhost:8080/api/mocks/{mock_id}/files/myfile.txt -O

# 删除文件
curl -X DELETE http://localhost:8080/api/mocks/{mock_id}/files/myfile.txt
```

## API 接口

### 创建 Mock API

**HTTP 示例：**
```http
POST /api/mocks
Content-Type: application/json

{
  "name": "测试接口",
  "port": 9090,
  "protocol": "http",
  "content": "Hello World",
  "charset": "UTF-8",
  "path": "/test",
  "method": "GET"
}
```

**HTTPS 示例：**
```http
POST /api/mocks
Content-Type: application/json

{
  "name": "HTTPS 测试接口",
  "port": 9443,
  "protocol": "https",
  "cert_file": "certs/server.crt",
  "key_file": "certs/server.key",
  "content": "Hello HTTPS",
  "charset": "UTF-8",
  "path": "/secure",
  "method": "GET"
}
```

**FTP 示例：**
```http
POST /api/mocks
Content-Type: application/json

{
  "name": "FTP 文件服务器",
  "port": 21,
  "protocol": "ftp",
  "ftp_mode": "passive",
  "ftp_root_dir": "./ftp_data/port_21",
  "ftp_user": "admin",
  "ftp_pass": "admin123",
  "ftp_passive_port_range": "50000-50100",
  "charset": "UTF-8"
}
```

### 获取所有 Mock API
```http
GET /api/mocks
```

### 获取单个 Mock API
```http
GET /api/mocks/:id
```

### 更新 Mock API
```http
PUT /api/mocks/:id
Content-Type: application/json

{
  "name": "更新后的名称",
  "content": "Updated content",
  "charset": "GBK"
}
```

### 删除 Mock API
```http
DELETE /api/mocks/:id
```

### FTP 文件管理 API

#### 列出文件
```http
GET /api/mocks/:id/files?path=/subfolder
```

#### 上传文件（限制 100MB）
```http
POST /api/mocks/:id/files
Content-Type: multipart/form-data

file: <binary file data>
path: /subfolder
```

#### 下载文件
```http
GET /api/mocks/:id/files/path/to/file.txt
```

#### 删除文件
```http
DELETE /api/mocks/:id/files/path/to/file.txt
```

## 项目结构

```
gomoco/
├── main.go                 # 主入口 (嵌入前端资源)
├── go.mod                  # Go 依赖管理
├── build.ps1               # Windows 构建脚本
├── build-linux.ps1         # Linux 交叉编译脚本
├── build-all.ps1           # 全平台构建脚本
├── config/                 # 配置文件目录
│   ├── mocks.yaml         # Mock API 配置 (自动生成)
│   └── mocks.yaml.example # 配置示例
├── internal/
│   ├── api/               # API 服务器
│   │   └── server.go      # 处理嵌入式静态资源
│   ├── models/            # 数据模型
│   │   └── mock.go
│   ├── server/            # Mock 服务器实现
│   │   ├── manager.go     # 服务器管理器
│   │   ├── http.go        # HTTP 服务器
│   │   └── tcp.go         # TCP 服务器
│   ├── storage/           # 持久化存储
│   │   └── storage.go     # YAML 文件存储
│   └── utils/             # 工具函数
│       └── charset.go     # 字符集转换
└── web/                   # 前端项目 (构建后嵌入到二进制)
    ├── package.json
    ├── vite.config.js
    ├── index.html
    └── src/
        ├── main.js
        ├── App.vue
        └── style.css
```

## 注意事项

- 同一端口只能被一个 Mock API 使用
- 删除 Mock API 会自动停止对应的服务并从配置文件中移除
- GBK 编码主要用于兼容老旧系统
- TCP Mock 会在接收到任何数据后立即返回配置的内容
- 所有配置自动保存到 `config/mocks.yaml`，重启后自动加载
- 首次运行会自动创建 `config` 目录
- **可执行文件是自包含的**，无需额外的前端文件或依赖
- 部署时只需要可执行文件和 `config` 目录（可选）

## 部署说明

### Linux 部署
```bash
# 1. 上传可执行文件
scp gomoco-linux-amd64 user@server:/opt/gomoco/

# 2. 添加执行权限
chmod +x /opt/gomoco/gomoco-linux-amd64

# 3. 运行（默认端口 8080）
./gomoco-linux-amd64

# 3a. 自定义端口运行
./gomoco-linux-amd64 -port 9000

# 4. 后台运行（可选）
nohup ./gomoco-linux-amd64 -port 8080 > gomoco.log 2>&1 &
```

### Docker 部署

项目已包含 `Dockerfile` 和 `docker-compose.yml`。

**快速部署：**
```bash
# 1. 构建 Linux 静态二进制
.\build-linux.ps1

# 2. 构建 Docker 镜像
docker build -t gomoco:latest .

# 3. 使用 Docker Compose 运行
docker-compose up -d

# 4. 查看日志
docker-compose logs -f

# 5. 访问
# http://localhost:8080
```

**优势：**
- 使用 `FROM scratch` 基础镜像，镜像极小
- 静态编译无需任何依赖
- 配置文件通过卷挂载持久化

## 后续计划

- [x] API 名称管理
- [x] 配置持久化 (YAML)
- [x] 单一可执行文件部署
- [x] 跨平台编译支持
- [ ] 支持动态响应（基于请求内容）
- [ ] 支持延迟响应
- [ ] 支持请求日志记录
- [ ] 支持响应模板
- [ ] 支持 HTTPS/TLS
- [ ] 导入/导出配置
- [ ] Docker 镜像

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！
