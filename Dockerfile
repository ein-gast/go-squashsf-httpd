FROM golang:1.25.1 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY internal ./internal
COPY *.go ./

RUN find . -type f
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o squashfs-httpd ./squashsf-httpd.go
