.PHONY: all test dockerbuild

TEMP_TAG=localhost/squashfs-httpd:latest
TARGET_BIN=squashfs-httpd

squashfs-httpd:
	CGO_ENABLED=0 go build -ldflags="-w -s" -o ${TARGET_BIN} ./cmd/squashsf-httpd

dockerbuild:
	docker build -t "${TEMP_TAG}" .
	ID=$$(docker create "${TEMP_TAG}") && \
	docker cp "$$ID:/app/squashfs-httpd" ${TARGET_BIN} && \
	docker container rm "$$ID"
	docker rmi --force "${TEMP_TAG}"

test:
	go test ./...

all: squashfs-httpd
