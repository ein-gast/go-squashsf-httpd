package server

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/ein-gast/go-squashsf-httpd/internal/filer"
	"github.com/ein-gast/go-squashsf-httpd/internal/logger"
	"github.com/ein-gast/go-squashsf-httpd/internal/settings"
)

type Server struct {
	srv *http.Server
	log *logger.Logger
	ctx context.Context
}

func NewServer(ctx context.Context, log *logger.Logger, cfg *settings.Settings) *Server {

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.BindAddr, cfg.BindPort),
		Handler: mux,
	}
	res := &Server{
		srv: srv,
		ctx: ctx,
	}

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
	resp.WriteHeader(404)
	resp.Header().Add("context-type", "text/plain; charset=utf-8")
	resp.Write([]byte("Not Found"))
}

func (srv *Server) archiveHandler(fs *filer.Filer, resp http.ResponseWriter, req *http.Request) {
	file, err := fs.PreOpen(req.RequestURI)
	if err != nil {
		srv.log.Msg(req.RemoteAddr, 404, req.Method, req.RequestURI, err.Error())
		resp.WriteHeader(404)
		resp.Header().Add("context-type", "text/plain; charset=utf-8")
		resp.Write([]byte("Not Found: "))
		resp.Write([]byte(err.Error()))
		return
	}
	defer file.Close()

	srv.log.Msg(req.RemoteAddr, 200, req.Method, req.RequestURI)
	resp.WriteHeader(200)
	resp.Header().Add("context-type", "text/plain; charset=utf-8")

	buf := make([]byte, 1024)
	_, err = io.CopyBuffer(resp, file, buf)
	if err != nil {
		srv.log.Msg(err.Error())
		resp.Write([]byte(err.Error()))
		return
	}
}
