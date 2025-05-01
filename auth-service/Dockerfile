# Этап 1: Кэширование зависимостей Go-модулей
FROM golang:1.24-alpine3.21 as modules

COPY go.mod go.sum /modules/

WORKDIR /modules

RUN go mod download

# Этап 2: Сборка приложения
FROM golang:1.24-alpine3.21 as builder

COPY --from=modules /go/pkg /go/pkg
COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/app ./cmd/app

# Этап 3: Создание итогового образа
FROM scratch

COPY --from=builder /app/config /config
COPY --from=builder /app/migrations /migrations
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/app"]