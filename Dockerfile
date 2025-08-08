FROM golang:1.24.4-alpine AS builder

RUN apk add --no-cache build-base git ca-certificates

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend ./

ENV CGO_ENABLED=1 GOOS=linux
RUN go build -ldflags="-s -w" -o /app/app ./cmd

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/app .

COPY backend/forum.db ./forum.db

COPY frontend/templates ./templates
COPY frontend/static ./static

COPY backend/migrations ./migrations

ENV SERVER_PORT=8080
ENV STATIC_DIR=/app/static
ENV TEMPLATE_DIR=/app/templates
ENV STATIC_PATH=/app/static
ENV TEMPLATE_PATH=/app/templates
ENV DATABASE_DRIVER=sqlite3
ENV DB_PATH=/app/forum.db
ENV PROVIDER=sqlite
ENV COOKIE_NAME=session_id

EXPOSE 8080

CMD ["./app"]