package server

import (
	"errors"
	"strings"
)

func pathUnderRoute(routePrefix string, urlPath string) (string, error) {
	suffix, found := strings.CutPrefix(urlPath, routePrefix)
	if !found {
		return "", errors.New("wrong prefix")
	}
	return "/" + strings.TrimLeft(suffix, "/"), nil
}
