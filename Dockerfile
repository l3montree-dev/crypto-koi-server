FROM golang:1.17 as build-env

WORKDIR /go/src/app

ENV GIN_MODE release

COPY . .

RUN go get -d -v ./...
RUN go build -o clodhopper-server ./cmd/clodhopper-server

FROM gcr.io/distroless/base

COPY --from=build-env /go/src/app/clodhopper-server /go/src/app/clodhopper-server
ENV GIN_MODE release
WORKDIR /go/src/app
EXPOSE 8080

CMD ["./clodhopper-server"]
