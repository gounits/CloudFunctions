FROM golang:alpine AS builder

WORKDIR /app
RUN apk update && apk add --no-cache upx ca-certificates
ENV GO111MODULE=on GOPROXY=https://goproxy.cn,direct
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod tidy
COPY . .
RUN go build -gcflags "all=-N -l" -tags timetzdata -o cf main.go
RUN upx cf

FROM busybox

WORKDIR /app

COPY --from=builder --chmod=777 /app/cf .
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/

EXPOSE 8080

ENTRYPOINT ["./cf"]