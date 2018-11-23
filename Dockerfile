FROM golang:1.10-alpine

RUN apk add --no-cache git

WORKDIR /go/src/airly-exporter
COPY *.go /go/src/airly-exporter/

RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /go/bin/
COPY --from=0 /go/bin/airly-exporter /go/bin/airly-exporter

EXPOSE 8080

# nobody
USER 65534

CMD ["./airly-exporter"]
