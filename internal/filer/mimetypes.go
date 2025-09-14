package filer

import (
	"path"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

type matchPair struct {
	t types.Type
	m func([]byte) bool
}

func AddMimeTypes() {
	moreTypes := []matchPair{
		{t: filetype.NewType("data-default", "application/octet-stream"), m: anyMatcher},
		{t: filetype.NewType("md", "text/plain"), m: anyMatcher},
		{t: filetype.NewType("htm", "text/html"), m: anyMatcher},
		{t: filetype.NewType("html", "text/html"), m: anyMatcher},
		{t: filetype.NewType("json", "application/json"), m: anyMatcher},
		{t: filetype.NewType("js", "text/javascript"), m: anyMatcher},
		{t: filetype.NewType("css", "text/css"), m: anyMatcher},
		{t: filetype.NewType("svg", "image/svg+xml"), m: anyMatcher},
	}

	for _, pair := range moreTypes {
		filetype.AddMatcher(pair.t, pair.m)
	}
}

func anyMatcher(buf []byte) bool {
	return true
}

func mimeByFilename(filePath string) types.MIME {
	ext := path.Ext(filePath)
	if len(ext) > 0 {
		ext = ext[1:]
	}

	if !filetype.IsSupported(ext) {
		return filetype.GetType("data-default").MIME
	}

	t := filetype.GetType(ext)
	return t.MIME
}
