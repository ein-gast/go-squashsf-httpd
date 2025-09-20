.PHONY: all test

squashfs-httpd:
	CGO_ENABLED=0 go build -ldflags="-w -s" -o squashfs-httpd ./squashsf-httpd.go

test:
	go test ./...

all: squashfs-httpd
