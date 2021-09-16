ARG GO_VERSION=1.16

# builder
FROM golang:${GO_VERSION} AS builder

RUN apt update

WORKDIR /api

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o app application.go

# Actual image
FROM ubuntu:20.04

LABEL maintainer="srihari.vishnu@gmail.com"
LABEL version="0.1"
LABEL description="This is the image for the shopify challenge server"

ARG DEBIAN_FRONTEND=noninteractive

RUN apt update && apt upgrade -y

WORKDIR /api
COPY --from=builder /api/app .

EXPOSE 5000

ENTRYPOINT ["./app"]