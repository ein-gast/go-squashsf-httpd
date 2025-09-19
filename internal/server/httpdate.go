package server

import "time"

func HttpDate(t time.Time) string {
	return t.Format(time.RFC1123)
}
