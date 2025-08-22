FROM golang:alpine AS builder

WORKDIR /app
RUN apk update && apk add --no-cache upx && rm -rf /var/cache/apk/*
ENV GO111MODULE=on GOPROXY=https://goproxy.cn,direct
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -ldflags "-s -w -X 'main.GO_VERSION=$(go version)' -X 'main.BUILD_TIME=`TZ=Asia/Shanghai date "+%F %T"`'" -o cf main.go
RUN upx cf

FROM busybox

COPY --from=builder --chmod=777 /app/cf .

EXPOSE 8080

ENTRYPOINT ["/cf"]