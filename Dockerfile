FROM golang:1.26-alpine AS builder

ENV GOPROXY="https://proxy.golang.org"
ENV GO111MODULE="on"
ENV NAT_ENV="production"
RUN apk add --no-cache git

WORKDIR /go/src/github.com/icco/numbers
COPY . .

EXPOSE 8080
ENV PORT=8080

RUN go build -o /go/bin/numbers .
CMD ["/go/bin/numbers"]
