FROM golang:1.13.4-alpine AS builder

ENV GO111MODULE=on

WORKDIR /repo

ADD go.mod .
ADD go.sum .

ADD . ./

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache bash tzdata

ENV TZ="Asia/Shanghai"

EXPOSE 6060

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -o api cmd/main.go

ENTRYPOINT ["/repo/api"]
