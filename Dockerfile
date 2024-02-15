FROM golang:1.22-alpine

RUN apk add --no-cache git

WORKDIR /app
COPY *.go .
COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go build

FROM alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=0 /app/airly-exporter /app/airly-exporter

EXPOSE 8080

# nobody
USER 65534

CMD ["./airly-exporter"]
