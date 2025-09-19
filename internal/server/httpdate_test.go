package server

import (
	"testing"
	"time"
)

func TestHttpDate(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"date 1", args{t: time.Date(2025, 12, 29, 16, 54, 11, 0, time.UTC)}, "Mon, 29 Dec 2025 16:54:11 UTC"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HttpDate(tt.args.t); got != tt.want {
				t.Errorf("HttpDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
