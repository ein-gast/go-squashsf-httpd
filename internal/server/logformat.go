package server

import (
	"net/http"
	"strings"
)

func logFormatDefault(code int, message string, req *http.Request) []any {
	addr := strings.Split(req.RemoteAddr, ":")
	return []any{addr[0], code, req.Method, req.RequestURI, message}
}
