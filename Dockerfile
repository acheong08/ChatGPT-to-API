# 设置基础镜像为 Golang 官方镜像
FROM golang:1.20.3-alpine

# 设置环境变量
# Reverse Proxy - Available on accessToken
# Default: https://bypass.churchless.tech/api/conversation
ENV API_REVERSE_PROXY 'https://bypass.churchless.tech/api/conversation'
ENV SERVER_HOST '0.0.0.0'

# 设置工作目录为 /app
WORKDIR /app

# 将本地应用程序复制到容器中
COPY . .

# 下载应用程序所需的依赖项
RUN go mod download

# 构建应用程序二进制文件
RUN go build -o app

# 暴露应用程序运行的端口
EXPOSE 8080

# 启动应用程序
CMD ["./app"]
