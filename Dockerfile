FROM golang:1.8

WORKDIR /go/src/airly-exporter
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8080

CMD ["airly-exporter"]
