FROM golang:1.14-alpine as builder

EXPOSE 8080

RUN apk add --no-cache --update git gcc musl-dev g++

WORKDIR /go/src/github.com/leboncoin/subot
COPY . .

WORKDIR /go/src/github.com/leboncoin/subot/services/analytics/cmd
RUN GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o analytics -ldflags "-linkmode external -extldflags -static"

FROM alpine:latest
RUN apk add --no-cache tzdata
WORKDIR /root
COPY --from=builder /go/src/github.com/leboncoin/subot/services/analytics/cmd/analytics .
RUN mkdir -p pkg/auth/dex
COPY --from=builder /go/src/github.com/leboncoin/subot/pkg/auth/dex/web ./pkg/auth/dex/web
COPY --from=builder /go/src/github.com/leboncoin/subot/config ./config
CMD ["./analytics"]
