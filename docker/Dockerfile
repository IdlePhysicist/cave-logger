FROM golang:1.13-alpine3.11 AS builder

RUN apk add \
  git \
  make \
  gcc \
  musl-dev \
  sqlite \
  sqlite-dev

RUN mkdir /root/.config

WORKDIR /go/src/github.com/IdlePhysicist/cave-logger

COPY go.* ./
#COPY go.sum .
RUN go mod download

COPY ./ ./

RUN make linux --always-make && make install

