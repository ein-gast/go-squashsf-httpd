package filer

import (
	"io"
	"io/fs"

	"github.com/ein-gast/go-squashsf-httpd/internal/settings"
	"github.com/h2non/filetype/types"
)

type Filer interface {
	Close()
	PreOpen(filePath string) (io.ReadCloser, fs.FileInfo, error)
	Mime(filePath string) types.MIME
}

func NewFiler(archive settings.ServedArchive) (Filer, error) {
	//filer, err := NewFilerCaleb(archive)
	filer, err := NewFilerDiskfs(archive)
	return filer, err
}
