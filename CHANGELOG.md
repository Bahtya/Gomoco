# Changelog

## [1.3.0] - 2025-12-15

### 新增功能
- ✨ **FTP 协议支持**: 支持 FTP Mock 服务器
- ✨ **FTP 主动/被动模式**: 支持 Active 和 Passive 模式
- ✨ **Web 文件管理**: 通过 Web API 管理 FTP 文件
- ✨ **文件上传限制**: 单个文件最大 100MB
- ✨ **FTP 认证**: 支持自定义用户名和密码

### 技术改进
- 🔧 使用 goftp/server 库实现 FTP 服务器
- 🔧 添加文件管理 API（列表、上传、下载、删除）
- 🔧 前端界面支持 FTP 配置
- 🔧 模型添加 FTP 相关字段
- 🔧 自动创建 FTP 根目录

## [1.2.0] - 2025-12-11

### 新增功能
- ✨ **HTTPS 协议支持**: 支持 HTTPS Mock API
- ✨ **SSL/TLS 证书**: 支持自定义证书和私钥文件
- ✨ **证书生成脚本**: 提供自签名证书生成工具

### 技术改进
- 🔧 HTTP 服务器支持 ListenAndServeTLS
- 🔧 前端界面添加证书文件配置
- 🔧 模型添加 CertFile 和 KeyFile 字段

## [1.1.0] - 2025-12-11

### 新增功能
- ✨ **API 名称字段**: 为每个 Mock API 添加描述性名称
- ✨ **配置持久化**: 使用 YAML 文件自动保存和加载配置
- ✨ **嵌入式前端**: 前端资源嵌入到可执行文件中
- ✨ **单一可执行文件部署**: 无需额外文件即可运行
- ✨ **跨平台编译**: 支持从 Windows 交叉编译 Linux 版本
- ✨ **多平台构建脚本**: 一键构建所有平台版本
- ✨ **静态编译**: Linux 版本完全静态链接，无运行时依赖
- ✨ **命令行参数**: 支持自定义端口和版本查看

### 技术改进
- 🔧 使用 Go embed.FS 嵌入静态资源
- 🔧 添加 storage 包处理 YAML 持久化
- 🔧 优化构建脚本，添加 `-ldflags="-s -w"` 减小文件大小
- 🔧 Linux 静态编译：`-extldflags '-static'` + `-tags netgo`
- 🔧 禁用 CGO (`CGO_ENABLED=0`) 确保完全静态
- 🔧 自动加载已保存的 Mock API 配置
- 🔧 使用 flag 包支持命令行参数解析

### 命令行参数
```bash
# 查看帮助
./gomoco -h

# 自定义端口
./gomoco -port 9000

# 查看版本
./gomoco -version
```

### 编译选项
```bash
# Linux 静态编译
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
go build -ldflags="-s -w -extldflags '-static'" -tags netgo -a

# 优势:
# - 无需 glibc 依赖
# - 可在任何 Linux 发行版运行
# - 适合 FROM scratch Docker 镜像
# - 避免库版本兼容问题
```

### 新增文件
- `internal/storage/storage.go` - 持久化存储层
- `config/mocks.yaml.example` - 配置文件示例
- `build-linux.ps1` - Linux 交叉编译脚本（静态链接）
- `build-all.ps1` - 全平台构建脚本（静态链接）
- `DEPLOYMENT.md` - 部署指南
- `Dockerfile` - Docker 镜像配置（FROM scratch）
- `docker-compose.yml` - Docker Compose 配置

### 变更
- 📝 更新 README.md，添加部署和构建说明
- 📝 更新 .gitignore，排除构建产物和配置文件
- 🔄 修改 API 模型，添加 Name 字段和 YAML 标签
- 🔄 更新前端界面，显示 API 名称

### 数据模型变更
```go
type MockAPI struct {
    ID       string // UUID
    Name     string // 新增: API 名称
    Port     int
    Protocol string
    Content  string
    Charset  string
    Path     string
    Method   string
    Status   string
}
```

## [1.0.0] - 2025-12-11

### 初始版本
- ✅ HTTP 和 TCP 协议支持
- ✅ UTF-8 和 GBK 字符集编码
- ✅ 固定报文内容响应
- ✅ 动态端口配置
- ✅ HTTP 路径和方法配置
- ✅ Vue 3 前端界面
- ✅ RESTful API 管理接口
