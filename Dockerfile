# Build stage - this is used and discarded
FROM golang:1.22-alpine AS build-env
WORKDIR /opt
COPY go.mod .
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o salad

# Deployable stage - final image is based on this
FROM alpine:latest
MAINTAINER jae@zhar.net
WORKDIR /opt
COPY --from=build-env /opt/salad .
ENV SERVER_ADDRESS=127.0.0.1:8888
CMD ["/opt/salad"]
