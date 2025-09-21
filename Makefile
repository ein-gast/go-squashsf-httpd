.PHONY: all test tag dockerbuild dist dist-linux-amd64 dist-linux-386 dist-linux-arm64

TEMP_TAG   = localhost/squashfs-httpd:latest
TARGET_BIN = squashfs-httpd
DIST_DIR=./var/dist/
BUILD_CMD  = CGO_ENABLED=0 go build -trimpath -ldflags="-w -s -X main.Version=$(TAG)" -o $(DIST_DIR)$(DIST_NAME)/${TARGET_BIN} ./cmd/squashsf-httpd
DIST_CMD   = cd $(DIST_DIR)$(DIST_NAME) && 7z a -tzip -so * > ../$(DIST_NAME).zip && cd .. && rm -rf $(DIST_NAME)

squashfs-httpd: DIST_DIR=.
squashfs-httpd: TAG=develop
squashfs-httpd: tag
	$(BUILD_CMD)

dockerbuild:
	docker build -t "${TEMP_TAG}" .
	ID=$$(docker create "${TEMP_TAG}") && \
	docker cp "$$ID:/app/squashfs-httpd" ${TARGET_BIN} && \
	docker container rm "$$ID"
	docker rmi --force "${TEMP_TAG}"

tag:
	git rev-parse --abbrev-ref HEAD > TAG
	grep -qF HEAD TAG && git tag --points-at HEAD > TAG || true
	grep -qF main TAG && echo -n "git-" > TAG && git rev-parse --short main >> TAG || true

test:
	go test ./...

dist-linux-amd64: DIST_NAME=linux-amd64-$(shell cat TAG)
dist-linux-amd64: TAG=$(shell cat TAG)
dist-linux-amd64: 
	mkdir -p $(DIST_DIR)$(DIST_NAME)
	GOOS=linux GOARCH=amd64 $(BUILD_CMD)
	$(DIST_CMD)

dist-linux-386: DIST_NAME=linux-i386-$(shell cat TAG)
dist-linux-386: TAG=$(shell cat TAG)
dist-linux-386:
	mkdir -p $(DIST_DIR)$(DIST_NAME)
	GOOS=linux GOARCH=386 $(BUILD_CMD)
	$(DIST_CMD)

dist-linux-arm64: DIST_NAME=linux-arm64-$(shell cat TAG)
dist-linux-arm64: TAG=$(shell cat TAG)
dist-linux-arm64:
	mkdir -p $(DIST_DIR)$(DIST_NAME)
	GOOS=linux GOARCH=arm64 $(BUILD_CMD)
	$(DIST_CMD)

dist: tag dist-linux-amd64 dist-linux-386 dist-linux-arm64

all: squashfs-httpd
