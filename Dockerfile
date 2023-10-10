#FROM golang:alpine AS builder
#WORKDIR /build
#ADD go.mod .
#COPY . .
#RUN go build -o service cmd/app/app.go
#FROM alpine
#WORKDIR /build
#COPY --from=builder /build/service /build/service
#EXPOSE 8080
#CMD [". /service"]

FROM golang:alpine AS builder

WORKDIR /usr/local/go/src/

ADD app/ /usr/local/go/src/

RUN go clean --modcache
RUN go build -mod=readonly -o app cmd/app/app.go

FROM alpine

COPY --from=builder /usr/local/go/src/app /
COPY --from=builder /usr/local/go/src/config.yml /

CMD ["/app"]