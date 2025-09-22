package filer

import (
	"os"
	"path"
	"testing"
)

func TestFilerDiskfs_Release(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Setup error 1: %v", err.Error())
	}
	file, err := os.Open(path.Join(cwd, "../../examples/data/potree-lion.sq"))
	if err != nil {
		t.Fatalf("Setup error 2: %v", err.Error())
	}
	fsobj, err := NewFilerDiskfs(file)
	if err != nil {
		t.Fatalf("Setup error 3: %v", err.Error())
	}

	r := t.Run("PreOpen sanity test", func(t *testing.T) {
		_, _, err := fsobj.PreOpen("/index.html")
		if err != nil {
			t.Fatalf("Unexpected error, PreOpen should've succeed: %v", err.Error())
		}
	})
	if !r {
		return
	}

	t.Run("Release reopens", func(t *testing.T) {
		file.Close()
		fsobj.Release() // must reopen file
		_, _, err := fsobj.PreOpen("/index.html")
		if err != nil {
			t.Fatalf("Unexpected error, PreOpen should've succeed: %v", err.Error())
		}
	})
}
