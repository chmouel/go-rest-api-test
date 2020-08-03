FROM golang AS builder
EXPOSE 8080

COPY . src
WORKDIR src

ENV GO111MODULE on
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w'  -o /usr/local/bin/go-rest-api-test

FROM alpine

COPY --from=builder /usr/local/bin/go-rest-api-test /usr/local/bin/go-rest-api-test

ENTRYPOINT ["/usr/local/bin/go-rest-api-test"]
