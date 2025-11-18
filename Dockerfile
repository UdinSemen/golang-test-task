FROM golang:1.25.1-alpine AS build
WORKDIR /app

ENV GOPATH=/go
ENV PATH=$PATH:$GOPATH/bin
ENV GO111MODULE=on
ENV GOOS=linux
ENV CGO_ENABLED=0

RUN apk add --no-cache git gcc libc-dev

COPY go.sum go.mod ./
RUN go mod download

COPY ./ ./

RUN go generate ./... \
    && go fmt ./... \
    && go build -ldflags="-s -w" -o go-test-service cmd/app/main.go

FROM alpine:latest

WORKDIR /app
RUN pwd
RUN apk add --no-cache tzdata ca-certificates
RUN cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime \
 && echo "Europe/Moscow" >  /etc/timezone \
 && echo "hosts: files dns" > /etc/nsswitch.conf

COPY --from=build /app/go-test-service ./go-test-service
COPY ./migrations ./migrations

RUN chmod +x ./go-test-service
ENTRYPOINT ["./go-test-service"]