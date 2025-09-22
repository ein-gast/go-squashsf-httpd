package filer

import (
	"os"
	"path"
	"testing"
)

func TestFilerDirDiskfs_Release(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Setup error 1: %v", err.Error())
	}
	fsobj, err := NewFilerDirDiskfs(path.Join(cwd, "../../examples/data"))
	if err != nil {
		t.Fatalf("Setup error 2: %v", err.Error())
	}

	r := t.Run("PreOpen sanity test", func(t *testing.T) {
		_, _, err := fsobj.PreOpen("/potree-lion.sq/index.html")
		if err != nil {
			t.Fatalf("Unexpected error, PreOpen should've succeed: %v", err.Error())
		}
	})
	if !r {
		return
	}

	t.Run("Release reopens", func(t *testing.T) {
		fsobj.Release() // must reopen file
		_, _, err := fsobj.PreOpen("/potree-lion.sq/index.html")
		if err != nil {
			t.Fatalf("Unexpected error, PreOpen should've succeed: %v", err.Error())
		}
	})
}
