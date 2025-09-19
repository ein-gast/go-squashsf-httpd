package server

import (
	"testing"

	"github.com/ein-gast/go-squashsf-httpd/internal/settings"
)

func Test_pathInArchive(t *testing.T) {
	type args struct {
		prefix  string
		urlPath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"A", args{"/", "/file"}, "/file", false},
		{"B", args{"/", "/file/path"}, "/file/path", false},
		{"C", args{"/file/", "/file/path"}, "/path", false},
		{"D", args{"/file", "/file/path/"}, "/path/", false},
		{"E", args{"/file", "/path/"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			route := settings.ServedArchive{
				UrlPrefix:   tt.args.prefix,
				ArchivePath: "~",
			}
			got, err := pathInArchive(route, tt.args.urlPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("pathInArchive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("pathInArchive() = %v, want %v", got, tt.want)
			}
		})
	}
}
