# =============================================================================
# API Driva - Dockerfile multi-stage
# =============================================================================
# DECISÃO: Build com Go 1.21, stage final Alpine para imagem pequena.
# =============================================================================

FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o api ./cmd/main.go

# Stage final
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/api .

EXPOSE 3000

CMD ["./api"]
