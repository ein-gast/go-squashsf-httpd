package filer

import (
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/diskfs/go-diskfs/backend"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/diskfs/go-diskfs/filesystem/squashfs"
	"github.com/ein-gast/go-squashsf-httpd/internal/settings"
	"github.com/h2non/filetype/types"
)

type fileStorage struct {
	file *os.File
}

func (stor *fileStorage) Sys() (*os.File, error) {
	return stor.file, nil
}
func (stor *fileStorage) Writable() (backend.WritableFile, error) {
	return stor.file, nil
}
func (stor *fileStorage) Close() error {
	return stor.file.Close()
}
func (stor *fileStorage) Read(b []byte) (int, error) {
	return stor.file.Read(b)
}
func (stor *fileStorage) ReadAt(p []byte, off int64) (n int, err error) {
	return stor.file.ReadAt(p, off)
}
func (stor *fileStorage) Seek(offset int64, whence int) (int64, error) {
	return stor.file.Seek(offset, whence)
}
func (stor *fileStorage) Stat() (os.FileInfo, error) {
	return stor.file.Stat()
}

type FilerDiskfs struct {
	archivePath string
	disk        *fileStorage
	diskStat    fs.FileInfo
	fs          filesystem.FileSystem
	// reader      squashfs.Reader
}

func NewFilerDiskfs(diskFile *os.File) (*FilerDiskfs, error) {
	var err error
	res := &FilerDiskfs{
		archivePath: diskFile.Name(),
	}
	res.diskStat, err = diskFile.Stat()
	if err != nil {
		return nil, err
	}
	res.disk = &fileStorage{
		file: diskFile,
	}
	res.fs, err = squashfs.Read(res.disk, res.diskStat.Size(), 0, 4096)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewFilerDiskfsFromRoute(archive settings.ServedArchive) (*FilerDiskfs, error) {
	res := &FilerDiskfs{
		archivePath: archive.ArchivePath,
	}
	file, err := os.Open(res.archivePath)
	if err != nil {
		return nil, err
	}
	return NewFilerDiskfs(file)
}

func (f *FilerDiskfs) Close() {
	f.fs.Close()
	f.disk.Close()
}

func (f *FilerDiskfs) Release() {
}

func (f *FilerDiskfs) PreOpen(filePath string) (io.ReadCloser, fs.FileInfo, error) {
	file, err := f.fs.OpenFile(filePath, os.O_RDONLY)
	if err != nil {
		return nil, nil, err
	}
	size, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, nil, err
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, nil, err
	}
	stat := &FileStat{
		fileName: path.Base(filePath),
		isDir:    false,
		size:     size,
		mTime:    f.diskStat.ModTime(),
	}
	return file, stat, nil
}

func (f *FilerDiskfs) Mime(filePath string) types.MIME {
	return mimeByFilename(filePath)
}
