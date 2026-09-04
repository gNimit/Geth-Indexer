# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /src
RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/geth-indexer ./cmd/geth-indexer

# Final runtime stage
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
COPY --from=builder /bin/geth-indexer /app/geth-indexer

USER appuser

EXPOSE 9090

ENTRYPOINT ["/app/geth-indexer"]
CMD ["-help"]
