FROM golang:1.17-alpine3.14 as builder

COPY . /go/src/certfetcher

WORKDIR /go/src/certfetcher

RUN CGO_ENABLED=0 GOOS=linux go build -v -o refresher ./cmd/refresher/...

FROM alpine:3.14
COPY --from=builder /go/src/certfetcher/refresher /bin/
ENTRYPOINT /bin/refresher