package settings

type Settings struct {
	BindAddr        string          // bind host
	BindPort        int             // bind port
	Archives        []ServedArchive // served archives
	DefaultChareset string          // text files charset (fddtfd to conten-type header)
	BufferSize      int             // read (from squashfs) buffer size
	ClientTimeout   float64         // client idle timeout
	AccessLog       string
	ErrorLog        string
}

type ServedArchive struct {
	ArchivePath string
	UrlPrefix   string
}

func NewSettings() *Settings {
	s := &Settings{
		BindAddr:        "127.0.0.1",
		BindPort:        8080,
		Archives:        make([]ServedArchive, 0, 1),
		DefaultChareset: "utf-8",
		BufferSize:      10240,
		ClientTimeout:   5.0,
		AccessLog:       "/dev/stderr",
		ErrorLog:        "/dev/stderr",
	}
	return s
}
