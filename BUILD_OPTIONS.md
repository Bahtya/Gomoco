# 构建选项说明

## 静态编译详解

### 为什么需要静态编译？

**问题：**
- 默认的 Go 程序依赖系统的 C 库（glibc）
- 不同 Linux 发行版的 glibc 版本不同
- 在 Ubuntu 编译的程序可能无法在 CentOS 上运行
- Alpine Linux 使用 musl libc 而非 glibc

**解决方案：**
静态编译将所有依赖打包到可执行文件中。

### 编译选项详解

#### 1. CGO_ENABLED=0
```bash
CGO_ENABLED=0
```
- **作用**: 禁用 CGO，不链接 C 库
- **优势**: 完全纯 Go 编译，无外部依赖
- **必需**: 静态编译的基础

#### 2. -ldflags="-s -w"
```bash
-ldflags="-s -w"
```
- **-s**: 去除符号表（symbol table）
- **-w**: 去除 DWARF 调试信息
- **效果**: 减小可执行文件大小 30-40%
- **缺点**: 无法使用 gdb 调试（生产环境不需要）

#### 3. -ldflags="-extldflags '-static'"
```bash
-ldflags="-extldflags '-static'"
```
- **作用**: 强制静态链接外部库
- **仅用于**: Linux 平台
- **配合**: CGO_ENABLED=0 使用

#### 4. -tags netgo
```bash
-tags netgo
```
- **作用**: 使用纯 Go 实现的网络库
- **替代**: 默认的 CGO 网络实现
- **优势**: 避免 DNS 解析依赖系统库

#### 5. -a
```bash
-a
```
- **作用**: 强制重新编译所有包
- **用途**: 确保所有包都使用静态编译选项
- **注意**: 编译时间会增加

### 完整命令示例

#### Windows 本地编译
```powershell
$env:CGO_ENABLED = "0"
go build -ldflags="-s -w" -tags netgo -a -o gomoco.exe main.go
```

#### Linux 静态编译（从 Windows）
```powershell
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -ldflags="-s -w -extldflags '-static'" -tags netgo -a -o gomoco-linux-amd64 main.go
```

#### macOS 编译（从 Windows）
```powershell
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -ldflags="-s -w" -tags netgo -a -o gomoco-darwin-amd64 main.go
```

### 验证静态编译

#### Linux 上验证
```bash
# 查看依赖库
ldd gomoco-linux-amd64

# 静态编译的输出应该是:
# not a dynamic executable

# 或者
file gomoco-linux-amd64
# 输出: statically linked
```

#### 测试兼容性
```bash
# 在不同发行版测试
docker run --rm -v $(pwd):/app alpine:latest /app/gomoco-linux-amd64
docker run --rm -v $(pwd):/app ubuntu:latest /app/gomoco-linux-amd64
docker run --rm -v $(pwd):/app centos:latest /app/gomoco-linux-amd64
```

### 文件大小对比

| 编译方式 | 文件大小 | 依赖 |
|---------|---------|------|
| 默认编译 | ~15MB | glibc |
| -ldflags="-s -w" | ~10MB | glibc |
| 完全静态 | ~10MB | 无 |

### Docker 镜像大小对比

| 基础镜像 | 镜像大小 | 说明 |
|---------|---------|------|
| FROM scratch | ~10MB | 仅可执行文件 |
| FROM alpine | ~15MB | 包含 shell 和工具 |
| FROM ubuntu | ~80MB | 完整 Linux 环境 |

### 性能影响

**启动时间：**
- 静态编译: 更快（无需加载动态库）
- 动态编译: 稍慢（需要加载 .so 文件）

**运行时性能：**
- 几乎无差异
- 网络库性能相同

**内存使用：**
- 静态编译: 略高（代码在内存中）
- 动态编译: 略低（共享库可被多进程共享）

### 最佳实践

#### 开发环境
```bash
# 快速编译，保留调试信息
go build -o gomoco main.go
```

#### 生产环境
```bash
# 完全优化的静态编译
CGO_ENABLED=0 go build -ldflags="-s -w -extldflags '-static'" -tags netgo -a -o gomoco main.go
```

#### 容器部署
```dockerfile
# 使用 scratch 基础镜像
FROM scratch
COPY gomoco /gomoco
CMD ["/gomoco"]
```

### 常见问题

#### Q: 为什么 macOS 不使用 -extldflags '-static'？
A: macOS 不支持完全静态链接，但 CGO_ENABLED=0 已经足够。

#### Q: 静态编译会影响性能吗？
A: 几乎没有影响，某些情况下甚至更快（无动态链接开销）。

#### Q: 可以在 Windows 上静态编译吗？
A: Windows 的 Go 程序默认就是静态编译的（除非使用 CGO）。

#### Q: 如何减小文件大小？
A: 使用 UPX 压缩（可选）：
```bash
upx --best --lzma gomoco-linux-amd64
# 可减小 50-70% 大小，但启动时需要解压
```

### 参考资源

- [Go 编译选项文档](https://pkg.go.dev/cmd/go)
- [静态编译最佳实践](https://github.com/golang/go/wiki/GoArm)
- [Docker 多阶段构建](https://docs.docker.com/build/building/multi-stage/)
