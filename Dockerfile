# 使用一个合适的基础镜像，例如官方的 Go 镜像
FROM golang:1.21.1

# 设置工作目录
WORKDIR /app

# 将 Go 源代码复制到工作目录
COPY . .

# 编译 Go 程序
RUN go build -o main .

# 指定容器启动时运行的命令
CMD ["./main"]
