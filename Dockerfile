#build stage
FROM golang:alpine AS builder
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on
RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk -u --no-cache add git
WORKDIR /go/src/app
COPY . .
RUN go build -o /go/bin/app -v cmd/producer/producer.go

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates mysql-client
COPY --from=builder /go/bin/app /dbsync/app
WORKDIR /dbsync
ENTRYPOINT /dbsync/app