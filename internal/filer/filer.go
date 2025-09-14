package filer

import (
	"os"
	"path"
	"strings"

	"github.com/CalebQ42/squashfs"
	"github.com/ein-gast/go-squashsf-httpd/internal/settings"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

type Filer struct {
	archivePath string
	arch        *os.File
	reader      squashfs.Reader
}

func NewFiler(archive settings.ServedArchive) (*Filer, error) {
	res := &Filer{
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

func (f *Filer) Close() {
	f.arch.Close()
}

func (f *Filer) PreOpen(filePath string) (*squashfs.File, error) {
	file, err := f.reader.OpenFile(f.normalizePath(filePath))
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f *Filer) Mime(file *squashfs.File, filePath string) types.MIME {
	ext := path.Ext(filePath)
	if len(ext) > 0 {
		ext = ext[1:]
	}

	if !filetype.IsSupported(ext) {
		return filetype.GetType("data").MIME
	}

	t := filetype.GetType(ext)
	return t.MIME
}

func (f *Filer) normalizePath(filePath string) string {
	return strings.TrimLeft(filePath, "/")
}
