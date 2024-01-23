# 使用一个基础镜像，例如官方的 Golang 镜像
FROM golang:1.21.1

# 设置工作目录
WORKDIR /backnas

# 复制 Go 代码到容器
COPY . .

# 构建 Go 应用程序
RUN GOPROXY=https://goproxy.io go build -o myapp cmd/ailbum.go
# 暴露应用程序的端口
EXPOSE 8001

# 定义容器启动命令
CMD ["./myapp"]
