# Gomoco Docker 镜像
# 使用静态编译的二进制文件，基于 scratch 最小镜像

FROM scratch

# 复制静态编译的可执行文件
COPY gomoco-linux-amd64 /gomoco

# 暴露 API 服务端口
EXPOSE 8080

# 运行应用（默认端口 8080）
CMD ["/gomoco"]

# 使用方法:
# 1. 构建 Linux 静态二进制: .\build-linux.ps1
# 2. 构建 Docker 镜像: docker build -t gomoco:latest .
# 3. 运行容器（默认端口）: 
#    docker run -d -p 8080:8080 -v $(pwd)/config:/config --name gomoco gomoco:latest
# 4. 运行容器（自定义端口）:
#    docker run -d -p 9000:9000 -v $(pwd)/config:/config --name gomoco gomoco:latest /gomoco -port 9000
#
# 注意: 配置文件会保存在容器内的 config 目录，建议挂载卷以持久化
