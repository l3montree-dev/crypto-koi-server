FROM golang:1.17 as build-env

WORKDIR /go/src/app

ENV GIN_MODE release

COPY . .

RUN go get -d -v ./...
RUN go build -o crypto-koi-api ./cmd/crypto-koi-api

FROM gcr.io/distroless/base

COPY --from=build-env /go/src/app/crypto-koi-api /go/src/app/crypto-koi-api
ENV GIN_MODE release
WORKDIR /go/src/app
EXPOSE 8080

CMD ["./crypto-koi-api"]
