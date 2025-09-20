package settings

type Settings struct {
	BindAddr        string  // bind host
	BindPort        int     // bind port
	DefaultChareset string  // text files charset (fddtfd to conten-type header)
	BufferSize      int     // read (from squashfs) buffer size
	ClientTimeout   float64 // client idle timeout
	AccessLog       string
	ErrorLog        string
	AccessLogOff    bool
	Archives        []ServedArchive    // served archives
	Directories     []ServedArchiveDir // served archive dirs
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
		BindAddr:        "127.0.0.1",
		BindPort:        8080,
		DefaultChareset: "utf-8",
		BufferSize:      10240,
		ClientTimeout:   5.0,
		AccessLog:       "/dev/stderr",
		ErrorLog:        "/dev/stderr",
		AccessLogOff:    false,
		Archives:        make([]ServedArchive, 0, 4),
		Directories:     make([]ServedArchiveDir, 0, 4),
	}
	return s
}
