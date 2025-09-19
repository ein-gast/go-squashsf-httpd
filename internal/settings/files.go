package settings

import (
	"fmt"
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

	fmt.Printf("[%s]\n", base)
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
	s.BindAddr = obj.BindAddr
	s.BindPort = obj.BindPort
	s.DefaultChareset = obj.DefaultChareset
	s.BufferSize = obj.BufferSize
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
