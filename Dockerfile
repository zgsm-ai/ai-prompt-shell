FROM golang:1.24.0 AS builder
WORKDIR /app
COPY . .

RUN go env -w CGO_ENABLED=0 && \
    go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct  

#
#go env -w GOPROXY=http://mirrors.sangfor.org/nexus/repository/go-proxy-group
#

RUN go mod tidy && go build -o ai-prompt-shell *.go

FROM alpine:3.21
#FROM golang:1.21-alpine
#FROM centos:7.6.1810
#时区设置
ENV env prod
ENV TZ Asia/Shanghai
WORKDIR /

COPY --from=builder /app/ai-prompt-shell /usr/local/bin
RUN chmod 755 /usr/local/bin/ai-prompt-shell
ENTRYPOINT ["/usr/local/bin/ai-prompt-shell"]

