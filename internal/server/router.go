package server

import (
	"errors"
	"strings"

	"github.com/ein-gast/go-squashsf-httpd/internal/settings"
)

func pathInArchive(route settings.ServedArchive, urlPath string) (string, error) {
	suffix, found := strings.CutPrefix(urlPath, route.UrlPrefix)
	if !found {
		return "", errors.New("wrong prefix")
	}
	return "/" + strings.TrimLeft(suffix, "/"), nil
}
