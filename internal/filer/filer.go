package filer

import (
	"io"
	"io/fs"
	"os"

	"github.com/ein-gast/go-squashsf-httpd/internal/settings"
	"github.com/h2non/filetype/types"
)

type Filer interface {
	PreOpen(filePath string) (io.ReadCloser, fs.FileInfo, error) // open file inside archive for reading
	Mime(filePath string) types.MIME                             // ask for mime type
	Close()                                                      // close archive
	Release()                                                    // release caches and buffers, close or reopen files
}

func NewFilerFromRoute(filefroute settings.ServedArchive) (Filer, error) {
	filer, err := NewFilerDiskfsFromRoute(filefroute)
	return filer, err
}

func NewFilerFromFd(diskFile *os.File) (Filer, error) {
	filer, err := NewFilerDiskfs(diskFile)
	return filer, err
}

func NewFilerDirFromRoute(dirroute settings.ServedArchiveDir) (Filer, error) {
	filer, err := NewFilerDirDiskfsFromRoute(dirroute)
	return filer, err
}
