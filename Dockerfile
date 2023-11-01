# 基础镜像
FROM golang:latest

# 设置工作目录
WORKDIR /app

# 添加空的go.mod文件
RUN touch go.mod

# 将代码复制到容器中
COPY . .

# 获取缺少的包
RUN go get golang.org/x/net/html/charset

# 编译Go程序
RUN go build -o main .

# 运行Go程序
CMD ["./main"]
