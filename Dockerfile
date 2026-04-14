FROM golang:1.26-alpine AS builder

ENV GOPROXY="https://proxy.golang.org"
ENV CGO_ENABLED=0

WORKDIR /app

# Cache dependency downloads separately from source changes.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o /numbers .

# ── Runtime image ─────────────────────────────────────────────────────────────
FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata

ENV NAT_ENV="production"
ENV PORT=8080

EXPOSE 8080

WORKDIR /app

# Run as a non-root user.
RUN adduser -S -u 1001 app
USER app

COPY --from=builder /numbers /app/numbers
# book.txt is read at runtime via os.ReadFile("book.txt") (relative path).
COPY book.txt /app/book.txt

CMD ["/app/numbers"]
