package filer

import (
	"io"
	"os"
	"strings"

	"github.com/CalebQ42/squashfs"
	"github.com/ein-gast/go-squashsf-httpd/internal/settings"
	"github.com/h2non/filetype/types"
)

type FilerCaleb struct {
	archivePath string
	arch        *os.File
	reader      squashfs.Reader
}

func NewFilerCaleb(archive settings.ServedArchive) (*FilerCaleb, error) {
	res := &FilerCaleb{
		archivePath: archive.ArchivePath,
	}
	var err error
	res.arch, err = os.Open(res.archivePath)
	if err != nil {
		return nil, err
	}
	res.reader, err = squashfs.NewReaderAtOffset(res.arch, 0)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (f *FilerCaleb) Close() {
	f.arch.Close()
}

func (f *FilerCaleb) PreOpen(filePath string) (io.ReadCloser, os.FileInfo, error) {
	file, err := f.reader.OpenFile(f.normalizePath(filePath))
	if err != nil {
		return nil, nil, err
	}
	stat, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}
	return file, stat, nil
}

func (f *FilerCaleb) Mime(filePath string) types.MIME {
	return mimeByFilename(filePath)
}

func (f *FilerCaleb) normalizePath(filePath string) string {
	return strings.TrimLeft(filePath, "/")
}
