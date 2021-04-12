FROM golang:1.14-alpine as builder

EXPOSE 8080

RUN apk add --no-cache git

WORKDIR /go/src/github.com/leboncoin/subot
COPY . .

WORKDIR /go/src/github.com/leboncoin/subot/services/replier/cmd
RUN GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o replier

FROM alpine:latest
RUN apk add --no-cache tzdata
WORKDIR /root
COPY --from=builder /go/src/github.com/leboncoin/subot/services/replier/cmd/replier .
COPY --from=builder /go/src/github.com/leboncoin/subot/config ./config
CMD ["./replier"]
