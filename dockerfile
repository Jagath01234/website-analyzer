FROM golang:1.21.7-alpine AS builder

ENV GOPROXY=https://proxy.golang.org,direct
ENV GO111MODULE=on

RUN apk update && apk add --no-cache git

WORKDIR /app

RUN git clone https://github.com/Jagath01234/website-analyzer.git

WORKDIR /app/website-analyzer

COPY config.json /app/website-analyzer/config.json

RUN go build -o website-analyzer .

FROM alpine:latest

COPY --from=builder /app/website-analyzer/website-analyzer /app/website-analyzer/website-analyzer
COPY --from=builder /app/website-analyzer/config.json /app/website-analyzer/config.json


WORKDIR /app/website-analyzer

EXPOSE 8080

CMD ["./website-analyzer"]
