FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Path: Dockerfile
FROM alpine

WORKDIR /app

COPY --from=builder /app/main .

ENTRYPOINT ["/app/main"]