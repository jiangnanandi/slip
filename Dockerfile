# 使用 Go 1.19 版本作为基础镜像
FROM golang:1.19 AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件,以便下载依赖
COPY go.mod go.sum ./

# 下载 Go 模块依赖
RUN go mod download

# 复制整个本地代码目录到容器中
COPY . .

# 构建 Go 应用程序
RUN go build -o slip main.go

# 使用轻量级的 Alpine 作为最终镜像
# FROM alpine:latest
FROM golang:1.19

# 创建工作目录
WORKDIR /app

# 从构建环境中复制可执行文件
COPY --from=builder /app/slip .

COPY config/config.yaml ./config/config.yaml

RUN mkdir -p /var/www/slip/notes && chmod 777 /var/www/slip/notes

# 暴露端口(如果您的应用需要)
EXPOSE 8084

# 设置入口点命令
CMD ["./slip"]
