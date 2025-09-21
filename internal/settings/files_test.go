package settings

import "testing"

func TestPathRelToAbs(t *testing.T) {
	type args struct {
		relPath  string
		basePath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"A", args{"/a/b/c", "/local"}, "/a/b/c"},
		{"B", args{"a/b/c", "/local"}, "/local/a/b/c"},
		{"C", args{"./a/b/c", "/local"}, "/local/a/b/c"},
		{"D", args{"", "/local"}, "/local"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PathRelToAbs(tt.args.relPath, tt.args.basePath); got != tt.want {
				t.Errorf("PathRelToAbs() = %v, want %v", got, tt.want)
			}
		})
	}
}
