package server

import (
	"testing"
)

func Test_pathUnderRoute(t *testing.T) {
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
		//
		{"X", args{"/dir", "/dir/file.sq/1"}, "/file.sq/1", false},
		{"Y", args{"/dir", "/dir/file.sq"}, "/file.sq", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pathUnderRoute(tt.args.prefix, tt.args.urlPath)
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
