package settings

import (
	"fmt"
	"runtime"

	"github.com/ein-gast/go-squashsf-httpd/internal/logger"
)

func PrintSetting(s Settings, ver string, log logger.Logger) {
	log.Msg(fmt.Sprintf("Version:\t%s", ver))
	log.Msg(fmt.Sprintf("Golang: \t%s", runtime.Version()))
	log.Msg(fmt.Sprintf("Listen: \t%s:%d", s.BindAddr, s.BindPort))
	log.Msg(fmt.Sprintf("Charset:\t%s", s.DefaultChareset))
	log.Msg(fmt.Sprintf("BuffSize:\t%d", s.BufferSize))
	log.Msg(fmt.Sprintf("C.Timeout:\t%f", s.ClientTimeout))
	if s.AccessLogOff {
		log.Msg(fmt.Sprintf("AccessLog:\t%s", "[disabled]"))
	} else {
		log.Msg(fmt.Sprintf("AccessLog:\t%s", s.AccessLog))
	}
	log.Msg(fmt.Sprintf("ErrorLog:\t%s", s.ErrorLog))
	if s.PidFileOff {
		log.Msg(fmt.Sprintf("PidFile:\t%s", "[disabled]"))
	} else {
		log.Msg(fmt.Sprintf("PidFile:\t%s", s.PidFile))
	}
	log.Msg("Routes:")
	for _, route := range s.Archives {
		log.Msg(fmt.Sprintf("file: %s -> %s", route.UrlPrefix, route.ArchivePath))
	}
	for _, route := range s.Directories {
		log.Msg(fmt.Sprintf(" dir: %s -> %s", route.UrlPrefix, route.DirectoryPath))
	}
}
