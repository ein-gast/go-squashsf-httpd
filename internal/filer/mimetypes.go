package filer

import (
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

type matchPair struct {
	t types.Type
	m func([]byte) bool
}

func AddMimeTypes() {
	moreTypes := []matchPair{
		{t: filetype.NewType("data", "application/octet-stream"), m: anyMatcher},
		{t: filetype.NewType("md", "text/plain"), m: anyMatcher},
	}

	for _, pair := range moreTypes {
		filetype.AddMatcher(pair.t, pair.m)
	}
}

func anyMatcher(buf []byte) bool {
	return true
}
