package filer

import (
	"io/fs"
	"time"
)

type FileStat struct {
	fileName string
	isDir    bool
	size     int64
	mTime    time.Time
}

func (info *FileStat) Name() string {
	return info.fileName
}
func (info *FileStat) Size() int64 {
	return info.size
}
func (info *FileStat) Mode() fs.FileMode {
	//fs.ModeDir
	//fs.ModePerm = 0777
	return 0
}
func (info *FileStat) ModTime() time.Time {
	return info.mTime
}
func (info *FileStat) IsDir() bool {
	return info.isDir
}
func (info *FileStat) Sys() any {
	return nil
}
