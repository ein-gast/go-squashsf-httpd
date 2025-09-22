package filer

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/ein-gast/go-squashsf-httpd/internal/settings"
	"github.com/h2non/filetype/types"
)

type FilerDirDiskfs struct {
	root       *os.Root
	rootPath   string
	opened     map[string]Filer
	openedLock *sync.RWMutex
}

func NewFilerDirDiskfs(dirPath string) (Filer, error) {
	res := &FilerDirDiskfs{
		root:       nil,
		rootPath:   dirPath,
		opened:     make(map[string]Filer, 8),
		openedLock: &sync.RWMutex{},
	}
	var err error
	res.root, err = res.openRoot()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewFilerDirDiskfsFromRoute(dirroute settings.ServedArchiveDir) (Filer, error) {
	return NewFilerDirDiskfs(dirroute.DirectoryPath)
}

func (f *FilerDirDiskfs) openRoot() (*os.Root, error) {
	root, err := os.OpenRoot(f.rootPath)
	if err != nil {
		return nil, err
	}
	return root, nil
}

func (f *FilerDirDiskfs) Close() {
	for _, fs := range f.opened {
		fs.Close()
	}
	f.root.Close()
}

func (f *FilerDirDiskfs) Release() {
	f.openedLock.Lock()
	defer f.openedLock.Unlock()
	for _, fs := range f.opened {
		fs.Close()
	}
	f.opened = make(map[string]Filer, 8)
	newRoot, err := f.openRoot()
	if err != nil {
		return
	}
	oldRoot := f.root
	f.root = newRoot
	oldRoot.Close()
}

func pathToParts(filePath string) []string {
	rest := "/" + strings.TrimLeft(filePath, "/")
	parts := []string{rest}
	for {
		rest = path.Dir(rest)
		parts = append(parts, rest)
		if rest == "/" {
			break
		}
	}
	return parts
}

func (f *FilerDirDiskfs) findSqFile(parts []string) (*os.File, string, error) {
	var file *os.File
	var err error
	var fullPath string
	for i := len(parts) - 1; i >= 0; i-- {
		fullPath = "./" + strings.TrimLeft(parts[i], "/")
		file, err = f.root.Open(fullPath)
		if err != nil {
			continue
		}
		stat, err := file.Stat()
		if err != nil || stat.IsDir() {
			continue
		}
		return file, parts[i], nil
	}
	return nil, "", errors.New("SQ file not found")
}

func (f *FilerDirDiskfs) getFiler(filePath string) (Filer, string, error) {
	parts := pathToParts(filePath)
	f.openedLock.RLock()
	for i := len(parts) - 1; i >= 0; i-- {
		filer, ok := f.opened[parts[i]]
		if ok {
			f.openedLock.RUnlock()
			// fs reused
			return filer, strings.TrimPrefix(filePath, parts[i]), nil
		}
	}
	f.openedLock.RUnlock()

	file, prefix, err := f.findSqFile(parts)
	if err != nil {
		return nil, "", err
	}

	filer, err := NewFilerFromFd(file)
	if err != nil {
		return nil, "", err
	}
	// fs opened
	f.openedLock.Lock()
	f.opened[prefix] = filer
	f.openedLock.Unlock()
	return filer, strings.TrimPrefix(filePath, prefix), nil
}

func (f *FilerDirDiskfs) PreOpen(filePath string) (io.ReadCloser, fs.FileInfo, error) {
	fs, innerPath, err := f.getFiler(filePath)
	if err != nil {
		return nil, nil, err
	}
	return fs.PreOpen(innerPath)
}

func (f *FilerDirDiskfs) Mime(filePath string) types.MIME {
	return mimeByFilename(filePath)
}
