package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ein-gast/go-squashsf-httpd/internal/filer"
	"github.com/ein-gast/go-squashsf-httpd/internal/logger"
	"github.com/ein-gast/go-squashsf-httpd/internal/settings"
)

type Server struct {
	srv     *http.Server
	log     *logger.Logger
	ctx     context.Context
	enc     string // default text encoding
	bufsize int    // copy buffer size
}

func NewServer(ctx context.Context, log *logger.Logger, cfg *settings.Settings) *Server {

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.BindAddr, cfg.BindPort),
		Handler: mux,
	}
	res := &Server{
		srv:     srv,
		ctx:     ctx,
		enc:     cfg.DefaultChareset,
		bufsize: cfg.BufferSize,
	}
	res.srv.SetKeepAlivesEnabled(true)

	rootHandled := false
	for _, archive := range cfg.Archives {
		if archive.UrlPrefix == "/" {
			rootHandled = true
		}
		handle, err := newPrefixHandler(res, archive)
		if err != nil {
			log.Msg("Cant start handling file:", archive.ArchivePath, "[", err.Error(), "]")
			continue
		}
		mux.HandleFunc(archive.UrlPrefix, handle)
	}
	if !rootHandled {
		mux.HandleFunc("/", res.nullHandler)
	}

	return res
}

func newPrefixHandler(srv *Server, archive settings.ServedArchive) (func(resp http.ResponseWriter, req *http.Request), error) {
	fs, err := filer.NewFiler(archive)
	if err != nil {
		return nil, err
	}
	return func(resp http.ResponseWriter, req *http.Request) {
		srv.archiveHandler(fs, resp, req)
	}, nil
}

func (srv *Server) Serve(log *logger.Logger) error {
	srv.log = log
	log.Msg("Serving...", srv.srv.Addr)
	err := srv.srv.ListenAndServe()
	log.Msg("Server stopped:", err.Error())
	return err
}

func (srv *Server) nullHandler(resp http.ResponseWriter, req *http.Request) {
	srv.log.Msg(req.RemoteAddr, 404, req.Method, req.RequestURI)
	resp.Header().Add("content-type", "text/plain; charset=utf-8")
	resp.Write([]byte("Not Found"))
	resp.WriteHeader(404)
}

func (srv *Server) archiveHandler(fs filer.Filer, resp http.ResponseWriter, req *http.Request) {
	filePath := req.URL.Path
	file, stat, err := fs.PreOpen(filePath)
	if err != nil {
		srv.log.Msg(req.RemoteAddr, 404, req.Method, req.RequestURI, err.Error())
		resp.Header().Add("content-type", "text/plain; charset=utf-8")
		resp.WriteHeader(404)
		resp.Write([]byte("Not Found: "))
		resp.Write([]byte(err.Error()))
		return
	}
	defer file.Close()

	var contentType string
	mime := fs.Mime(filePath)
	if mime.Type == "text" {
		contentType = mime.Value + "; charset=" + srv.enc
	} else {
		contentType = mime.Value
	}

	srv.log.Msg(req.RemoteAddr, 200, contentType, req.Method, req.RequestURI)
	resp.Header().Add("content-type", contentType)
	resp.Header().Add("content-length", fmt.Sprintf("%d", stat.Size()))
	resp.Header().Add("last-modified", HttpDate(stat.ModTime()))
	resp.Header().Add("x-name", stat.Name())

	if !isModifiedSince(req.Header.Get("if-modified-since"), stat.ModTime()) {
		resp.WriteHeader(304) // not modified
		return
	}

	resp.WriteHeader(200)

	if req.Method == http.MethodHead {
		return
	}

	buf := make([]byte, srv.bufsize)
	_, err = io.CopyBuffer(resp, file, buf)
	if err != nil {
		srv.log.Msg(err.Error())
		resp.Write([]byte(err.Error()))
		return
	}
}

func isTimeEqualSoft(a, b time.Time) bool {
	sub := a.Sub(b)
	if sub < 0 && sub > -time.Second {
		return true
	}
	if sub > 0 && sub < time.Second {
		return true
	}
	return false
}

func isModifiedSince(headerTime string, mtime time.Time) bool {
	if len(headerTime) == 0 {
		return true
	}
	htime, err := time.Parse(time.RFC1123, headerTime)
	if err != nil {
		return true
	}
	if htime.After(mtime) || isTimeEqualSoft(htime, mtime) {
		return false
	}
	return true
}
