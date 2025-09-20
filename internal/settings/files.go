package settings

import (
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"
)

type YamlRoute struct {
	Prefix string `yaml:"prefix"`
	Squash string `yaml:"squash"`
}

type YamlSettings struct {
	BindAddr        string      `yaml:"bind_addr"`
	BindPort        int         `yaml:"bind_port"`
	DefaultChareset string      `yaml:"charset"`
	BufferSize      int         `yaml:"buffer"`
	ClientTimeout   float64     `yaml:"client_timeout"`
	AccessLog       string      `yaml:"access_log"`
	ErrorLog        string      `yaml:"error_log"`
	Routes          []YamlRoute `yaml:"routes"`
}

func Load(cfgPath string) (*Settings, error) {
	base := path.Dir(cfgPath)
	if !path.IsAbs(base) {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		base = path.Join(cwd, base)
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	y := &YamlSettings{}
	err = yaml.Unmarshal(data, y)
	if err != nil {
		return nil, err
	}
	// transrorming log paths
	if len(y.ErrorLog) > 0 && !path.IsAbs(y.ErrorLog) {
		y.ErrorLog = path.Join(base, y.ErrorLog)
	}
	if len(y.ErrorLog) > 0 && !path.IsAbs(y.AccessLog) {
		y.AccessLog = path.Join(base, y.AccessLog)
	}
	// transrorming route paths
	for i := range y.Routes {
		sq := y.Routes[i].Squash
		if !path.IsAbs(sq) {
			y.Routes[i].Squash = path.Join(base, sq)
		}
	}

	return y.ToSetting(), nil
}

func (obj *YamlSettings) ToSetting() *Settings {
	s := NewSettings()
	s.BindAddr = strDefault(obj.BindAddr, s.BindAddr)
	s.BindPort = intDefault(obj.BindPort, s.BindPort)
	s.DefaultChareset = strDefault(obj.DefaultChareset, s.DefaultChareset)
	s.BufferSize = intDefault(obj.BufferSize, s.BufferSize)
	s.ClientTimeout = obj.ClientTimeout
	s.AccessLog = strDefault(obj.AccessLog, s.AccessLog)
	s.ErrorLog = strDefault(obj.ErrorLog, s.ErrorLog)
	s.Archives = make([]ServedArchive, 0, len(obj.Routes))
	for _, r := range obj.Routes {
		s.Archives = append(s.Archives, ServedArchive{
			UrlPrefix:   r.Prefix,
			ArchivePath: r.Squash,
		})
	}
	return s
}

func (s *Settings) ToYaml() *YamlSettings {
	obj := &YamlSettings{
		BindAddr:        s.BindAddr,
		BindPort:        s.BindPort,
		DefaultChareset: s.DefaultChareset,
		BufferSize:      s.BufferSize,
		ClientTimeout:   s.ClientTimeout,
		AccessLog:       s.AccessLog,
		ErrorLog:        s.ErrorLog,
		Routes:          make([]YamlRoute, 0, len(s.Archives)),
	}
	for _, r := range s.Archives {
		obj.Routes = append(obj.Routes, YamlRoute{
			Prefix: r.UrlPrefix,
			Squash: r.ArchivePath,
		})
	}
	return obj
}

func strDefault(val string, def string) string {
	if len(val) > 0 {
		return val
	}
	return def
}

func intDefault(val int, def int) int {
	if val != 0 {
		return val
	}
	return def
}
