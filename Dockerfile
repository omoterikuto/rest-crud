FROM golang:1.16

ENV GO111MODULE=on

RUN mkdir /go/src/app

WORKDIR /go/src

COPY go.mod go.sum ./

RUN go mod download

WORKDIR /go/src/app

ADD . /go/src/app

RUN go get github.com/go-sql-driver/mysql

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o /main ./cmd

FROM alpine:3.12

COPY --from=builder /main .

ENV PORT=${PORT}

ENTRYPOINT ["/main web"]
