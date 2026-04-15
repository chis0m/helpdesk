# Multi-stage image for HelpDesk backend. 

# Stage 1: Build
FROM golang:1.26.2-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build \
    -ldflags="-w -s" \
    -o /app/bin/helpdesk \
    ./cmd/server/main.go

# Stage 2: Runtime
FROM alpine:3.21

RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    wget \
    && update-ca-certificates

RUN addgroup -g 1000 helpdesk && \
    adduser -D -u 1000 -G helpdesk helpdesk

WORKDIR /app

COPY --from=builder /app/bin/helpdesk /app/helpdesk
COPY --from=builder /app/migrations /app/migrations
# Default vuln image; secure: docker build --build-arg ENV_VARIANT=secure ...
ARG ENV_VARIANT=vuln
COPY .env.${ENV_VARIANT}.local .env

ENV PORT=8080

RUN chown -R helpdesk:helpdesk /app

USER helpdesk

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
    CMD /bin/sh -c 'wget --no-verbose --tries=1 --spider "http://127.0.0.1:${PORT:-8080}/healthz" || exit 1'

CMD ["/app/helpdesk"]
