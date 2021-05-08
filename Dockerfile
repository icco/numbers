FROM golang:1.16-alpine as builder

ENV GOPROXY="https://proxy.golang.org"
ENV GO111MODULE="on"
ENV NAT_ENV="production"
RUN apk add --no-cache git

WORKDIR /go/src/github.com/icco/numbers
COPY . .

RUN go build -o /go/bin/numbers .
CMD /go/bin/numbers
