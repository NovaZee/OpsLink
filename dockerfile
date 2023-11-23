# 基于 golang 官方镜像作为基础镜像
FROM golang:latest

# 设置工作目录
WORKDIR /app

# 将本地的代码复制到镜像中的工作目录
COPY . .

# 构建 Go 程序
RUN go build -o main .

# 暴露程序运行所需的端口
EXPOSE 8080

# 运行 Go 程序
CMD ["./main"]