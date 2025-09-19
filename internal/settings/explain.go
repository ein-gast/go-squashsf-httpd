package settings

import (
	"fmt"

	"github.com/ein-gast/go-squashsf-httpd/internal/logger"
)

func PrintSetting(s Settings, log *logger.Logger) {
	log.Msg(fmt.Sprintf("Listen:\t%s:%d", s.BindAddr, s.BindPort))
	log.Msg(fmt.Sprintf("Charset:\t%s", s.DefaultChareset))
	log.Msg(fmt.Sprintf("BuffSize:\t%d", s.BufferSize))
	log.Msg("Routing:")
	for _, route := range s.Archives {
		log.Msg(fmt.Sprintf("%s -> %s", route.UrlPrefix, route.ArchivePath))
	}
}
