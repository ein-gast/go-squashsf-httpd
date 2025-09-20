FROM golang:1.25.1 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal

RUN find . -type f
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o squashfs-httpd ./cmd/squashsf-httpd
