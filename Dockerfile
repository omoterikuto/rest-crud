FROM golang:1.16

ENV GO111MODULE=on

RUN mkdir /go/src/app

WORKDIR /go/src

COPY go.mod go.sum ./

RUN go mod download

WORKDIR /go/src/app

ADD . /go/src/app

RUN go get github.com/go-sql-driver/mysql
