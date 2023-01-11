FROM devopsworks/golang-upx:1.18 AS builder

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
ARG user
ENV HOST_USER=$user

WORKDIR /repo

# all the steps are cached
ADD go.mod .
ADD go.sum .
# if go.mod/go.sum not changed, this step is also cached
RUN go mod download

ADD . ./
RUN go mod vendor

RUN export GDD_VER=$(go list -mod=vendor -m -f '{{ .Version }}' github.com/unionj-cloud/go-doudou/v2) && \
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags="-s -w -X 'github.com/unionj-cloud/go-doudou/v2/framework/buildinfo.BuildUser=$HOST_USER' -X 'github.com/unionj-cloud/go-doudou/v2/framework/buildinfo.BuildTime=$(date)' -X 'github.com/unionj-cloud/go-doudou/v2/framework/buildinfo.GddVer=$GDD_VER'" -mod vendor -o api cmd/main.go && \
strip api && /usr/local/bin/upx api

FROM alpine:3.14

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ="Asia/Shanghai"

WORKDIR /repo

COPY --from=builder /repo/api ./

COPY .env* ./

ENTRYPOINT ["/repo/api"]
