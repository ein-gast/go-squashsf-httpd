# squashsf-httpd

This software is a HTTP-server which gets static files directly from SquashFS archives without mounting it.

The server is written in golang. [github.com/diskfs/go-diskfs] is used to access SquashFS.

# Building from surce

Building requires golang>=1.24.
```bash
git clone https://github.com/ein-gast/go-squashsf-httpd.git
make all
stat squashfs-httpd
```

You also can build using docker (golang is not required on host).
```bash
git clone https://github.com/ein-gast/go-squashsf-httpd.git
make dockerbuild
stat squashfs-httpd.bin
```

## Installation and configuration

If you used `make all` or downloaded release then you may place `squashfs-httpd` binary wherever you want.

Alternatively you can use golang package manager.
```bash
go install github.com/ein-gast/go-squashsf-httpd/cmd/squashsf-httpd@latest
```

The easiest way to serve files from SquashFS is:
```bash
./squashfs-httpd -host 127.0.0.1 -port 8080 -squash ./examples/data/potree-lion.sq
```

This command provides an example of point cloud at url `http://127.0.0.1:8080/index.html`. This example is borrowed from [Potree](https://github.com/potree/potree) project.

For more details see: `./squashfs-httpd --help`

A configuation file is required to serve more than one SquashFS file.
```bash
./squashfs-httpd -config squashfs-httpd.yaml
```

Congiguration file eample.
```yaml
#  -- squashfs-httpd.yaml --
bind_addr: 127.0.0.1
bind_port: 8080
charset: utf-8
buffer: 10240
error_log: "./var/logs/error.log"
access_log: "./var/logs/access.log"
access_log_off: false
client_timeout: 5.0
routes:
  - prefix: /one/
    squash: ./examples/data/potree-lion.sq
  - prefix: /two/
    squashdir: ./examples/data/
```

If **USR1** signal is got then the server reopens logs by. If **USR2** signal is got then the server closes files which are opened below `squashdir` routes.

## Usage pattern

**squashsf-httpd** is designed to work in containers together with nginx and serve large number of small files packed in SquashFS.

The data which led to the creation of this software were point clouds and tile caches. These data are usually "tiled" once and then stored as millions of nested folders and files. These folders are hard to back up or move. A nice way to solve this inconvenience is to pack the files into SquashFS and mount it.

Mounting SquashFS starts being a problem when new data emerges on a daily basis. Someone should manage access rights to let applications mount new files. Another problem is an application in an unprivileged container; such applications cannot use even FUSE without tricks. So, it is convenient to have a small HTTP-server which could get files from SquashFS directly, and this is the main use-case for **squashfs-httpd**.

# Links
- https://github.com/plougher/squashfs-tools
- https://github.com/diskfs/go-diskfs
- https://github.com/CalebQ42/squashfs
- https://github.com/h2non/filetype
- https://github.com/potree/potree

