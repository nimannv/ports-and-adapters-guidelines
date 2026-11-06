# Stage 1: Build binaries
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy dependency manifests
COPY go.mod go.sum* ./
RUN go mod download || true

# Copy source code
COPY . .

# Run tidy and build statically linked binaries
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/api ./bootstrap/api && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/task-cli ./bootstrap/cli

# Stage 2: Minimal runtime image
FROM alpine:3.21

RUN apk --no-cache add ca-certificates curl tzdata

WORKDIR /app

COPY --from=builder /bin/api /app/api
COPY --from=builder /bin/task-cli /usr/local/bin/task-cli
COPY configs /app/configs

EXPOSE 8000

ENV CONFIGS_DIR=/app/configs \
    API_HOST=0.0.0.0 \
    API_PORT=8000

ENTRYPOINT ["/app/api"]

