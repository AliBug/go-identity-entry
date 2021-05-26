FROM golang:1.16.4-alpine as builder

WORKDIR /usr/src/app

# 这里将Golang依赖定义相关文件的copy放到最前面
COPY go.mod go.sum ./app/main.go ./

ENV GOPROXY=https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
  apk add --no-cache upx ca-certificates tzdata

RUN go env && go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o server

FROM alpine:3.13 
LABEL MAINTAINER="leijinchao@gmail.com"

# 🍉 此句似乎无用 RUN apk --no-cache add ca-certificates
# 🍎 Workdir 似乎应该修改
WORKDIR /app/
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/src/app/server ./

EXPOSE 8080

ENTRYPOINT ./server