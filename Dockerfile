FROM golang:alpine as builder
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app cmd/app/app.go
# Финальный этап, копируем собранное приложение
FROM alpine
COPY config.yml /config.yml
COPY --from=builder app /bin/app
EXPOSE 8080
ENTRYPOINT ["/bin/app"]