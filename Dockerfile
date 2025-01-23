# stage 1 - build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# download modules
COPY go.mod go.sum ./
RUN go mod download

# copy project
COPY . .

# build
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server


# stage 2 - run
FROM alpine:latest

# sops 설치
ARG SOPS_VERSION=3.9.3
RUN apk add --no-cache curl bash && \
    curl -L "https://github.com/getsops/sops/releases/download/v${SOPS_VERSION}/sops-v${SOPS_VERSION}.linux.amd64" \
    -o /usr/local/bin/sops && \
    chmod +x /usr/local/bin/sops

WORKDIR /app

# copy the binary and config directory
COPY --from=builder /app/main .
COPY --from=builder /app/config/encrypted-property.yml /app/config/
COPY --from=builder /app/internal/template/templates /app/internal/template/templates

# sops 복호화, start main
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]
