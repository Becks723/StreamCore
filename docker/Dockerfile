# 构建阶段
FROM golang:1.25 AS builder

# 定义构建参数，构建镜像时需要传入 SERVICE 参数来表示构建的模块
ARG SERVICE

# 设置环境变量
ENV TZ=Asia/Shanghai
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# 创建工作目录
RUN mkdir -p /app
RUN mkdir -p /app/output
WORKDIR /app

# 拷贝所有文件到工作目录
COPY . .
# go mod依赖
RUN go mod tidy

# 编译
RUN bash ./docker/script/build.sh ${SERVICE}

# 打包阶段
FROM alpine

ARG SERVICE

# 设置环境变量
ENV TZ=Asia/Shanghai
ENV SERVICE=${SERVICE}

# 工作目录
WORKDIR /app

# 拷贝bin文件和启动脚本，设置权限
COPY --from=builder /app/output /app/output
COPY --from=builder /app/docker/script/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# 需要注意的是，这个命令是在程序执行时运行，不能直接在这里填入 SERVICE 参，需要绕个弯设置为环境变量
CMD ["sh", "-c", "./entrypoint.sh"]
