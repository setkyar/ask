FROM golang:latest

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

COPY . .

RUN go mod tidy

RUN go build -o /go/bin/app

ENTRYPOINT ["/go/bin/app"]

CMD ["--help"]