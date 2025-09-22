package settings

import (
	"fmt"
	"os"
	"path"
)

type Settings struct {
	BindAddr         string             // bind host
	BindPort         int                // bind port
	DefaultChareset  string             // text files charset (fddtfd to conten-type header)
	BufferSize       int                // read (from squashfs) buffer size
	ClientTimeout    float64            // client idle timeout
	AccessLog        string             // path to access log
	ErrorLog         string             // path to error log
	AccessLogOff     bool               // true = do not write to access log
	PidFile          string             // path to pid file
	PidFileOff       bool               // true = do not create or use pid files
	DataCacheOff     bool               // true = do not use data cache
	DataCacheCount   int                // max cache entries count
	DataCacheEntSize int                // max cache entry size, larger files are not cached
	Archives         []ServedArchive    // served archives
	Directories      []ServedArchiveDir // served archive dirs
}

type ServedArchive struct {
	ArchivePath string
	UrlPrefix   string
}

type ServedArchiveDir struct {
	DirectoryPath string
	UrlPrefix     string
}

func NewSettings() *Settings {
	s := &Settings{
		BindAddr:         "127.0.0.1",
		BindPort:         8080,
		DefaultChareset:  "utf-8",
		BufferSize:       10240,
		ClientTimeout:    5.0,
		AccessLog:        "/dev/stderr",
		ErrorLog:         "/dev/stderr",
		AccessLogOff:     false,
		PidFile:          path.Join(defaultPidFolder(), "squashfs-httpd.pid"),
		PidFileOff:       false,
		DataCacheOff:     false,
		DataCacheCount:   500,
		DataCacheEntSize: 1024 * 1024,
		Archives:         make([]ServedArchive, 0, 4),
		Directories:      make([]ServedArchiveDir, 0, 4),
	}
	return s
}

func defaultPidFolder() string {
	runDir := "/run"
	uid := os.Getuid()
	isRoot := uid == 0
	hasRun := false
	stat, err := os.Stat(runDir)
	if err == nil && stat.IsDir() {
		hasRun = true
	}
	var pidFolder string
	if hasRun {
		if isRoot {
			pidFolder = runDir
		} else {
			pidFolder = path.Join(runDir, "user", fmt.Sprintf("%d", uid))
		}
	} else {
		pidFolder = os.TempDir()
	}
	return pidFolder
}
