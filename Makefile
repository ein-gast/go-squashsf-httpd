.PHONY: all

squashfs-httpd:
	CGO_ENABLED=0 go build -ldflags="-w -s" -o squashfs-httpd ./squashsf-httpd.go

all: squashfs-httpd
