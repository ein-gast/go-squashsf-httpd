package settings

type Settings struct {
	BindAddr        string
	BindPort        int
	Archives        []ServedArchive
	DefaultChareset string
	BufferSize      int
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
	}
	return s
}
