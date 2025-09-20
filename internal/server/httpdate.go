package server

import "time"

func HttpDate(t time.Time) string {
	return t.Format(time.RFC1123)
}

func IsTimeEqualSoft(a, b time.Time) bool {
	sub := a.Sub(b)
	if sub < 0 && sub > -time.Second {
		return true
	}
	if sub > 0 && sub < time.Second {
		return true
	}
	return false
}

func IsModifiedSince(headerTime string, mtime time.Time) bool {
	if len(headerTime) == 0 {
		return true
	}
	htime, err := time.Parse(time.RFC1123, headerTime)
	if err != nil {
		return true
	}
	if htime.After(mtime) || IsTimeEqualSoft(htime, mtime) {
		return false
	}
	return true
}
