FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN go build -o service cmd/app/app.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/service /build/service

EXPOSE 8080

CMD [". /service"]