# --- 第一阶段：编译阶段 ---
FROM golang:1.24.0-alpine AS builder

# 设置环境变量
# CGO_ENABLED=0 确保静态编译，这样在运行阶段不需要 c 库
# GOOS=linux 目标操作系统为 linux
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 设置工作目录
WORKDIR /build

# 先拷贝依赖文件，利用 Docker 缓存
COPY go.mod ./
RUN go mod download

# 拷贝源代码
COPY . .

# 编译成可执行文件，命名为 app
RUN go build -ldflags="-s -w" -o kube-vue-admin .

# --- 第二阶段：运行阶段 ---
FROM alpine:latest

# 安装基础库，如时区数据和证书（如果需要请求 https）
RUN apk --no-cache add ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

# 从编译阶段拷贝生成的可执行文件
COPY --from=builder /build/kube-vue-admin .
COPY --from=builder /build/configs/ ./configs/

# 暴露端口（根据你的项目修改）
EXPOSE 9000

# 启动程序
ENTRYPOINT ["./kube-vue-admin"]

# cd kube-vue-admin
# cp example/backend/Dockerfile ./
# docker build -t kube-vue-admin:v1.0 .