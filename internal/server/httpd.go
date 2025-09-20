package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/ein-gast/go-squashsf-httpd/internal/filer"
	"github.com/ein-gast/go-squashsf-httpd/internal/logger"
	"github.com/ein-gast/go-squashsf-httpd/internal/pool"
	"github.com/ein-gast/go-squashsf-httpd/internal/settings"
)

type Server struct {
	srv   *http.Server
	elog  logger.Logger
	alog  logger.Logger
	ctx   context.Context
	enc   string // default text encoding
	bsize int
	bpool *pool.BufferPool
}

func NewServer(
	ctx context.Context,
	elog logger.Logger,
	alog logger.Logger,
	cfg *settings.Settings,
) *Server {
	srv := &http.Server{
		Addr:        fmt.Sprintf("%s:%d", cfg.BindAddr, cfg.BindPort),
		BaseContext: func(net.Listener) context.Context { return ctx },
	}
	res := &Server{
		srv:  srv,
		ctx:  ctx,
		elog: elog,
		alog: alog,
	}
	res.srv.SetKeepAlivesEnabled(true)

	res.ApplyConfig(cfg)
	return res
}

func (srv *Server) ApplyConfig(cfg *settings.Settings) {
	srv.enc = cfg.DefaultChareset
	srv.bsize = cfg.BufferSize
	srv.bpool = pool.NewBufferPool(cfg.BufferSize)

	mux := http.NewServeMux()
	rootHandled := false
	for _, archive := range cfg.Archives {
		if archive.UrlPrefix == "/" {
			rootHandled = true
		}
		handle, err := newPrefixHandler(srv, archive)
		if err != nil {
			srv.elog.Msg("Cant start handling file:", archive.ArchivePath, "[", err.Error(), "]")
			continue
		}
		mux.HandleFunc(archive.UrlPrefix, handle)
	}
	if !rootHandled {
		mux.HandleFunc("/", srv.nullHandler)
	}
	srv.srv.Handler = mux

	timeout := time.Second * time.Duration(cfg.ClientTimeout)
	srv.srv.ReadTimeout = timeout
	srv.srv.WriteTimeout = timeout
}

func newPrefixHandler(srv *Server, archive settings.ServedArchive) (func(resp http.ResponseWriter, req *http.Request), error) {
	fs, err := filer.NewFiler(archive)
	if err != nil {
		return nil, err
	}
	return func(resp http.ResponseWriter, req *http.Request) {
		srv.archiveHandler(fs, archive, resp, req)
	}, nil
}

func (srv *Server) Serve() error {
	srv.elog.Msg("Serving...", srv.srv.Addr)
	err := srv.srv.ListenAndServe()
	srv.elog.Msg("Server stopped:", err.Error())
	return err
}
func (srv *Server) Shutdown() error {
	srv.elog.Msg("Shutting down server...")
	err := srv.srv.Shutdown(context.Background())
	if err != nil {
		srv.elog.Msg("Shutdown error:", err.Error())
	}
	return err
}

func (srv *Server) nullHandler(resp http.ResponseWriter, req *http.Request) {
	srv.writeError(404, "Not Found", resp, req)
}

func (srv *Server) writeError(code int, message string, resp http.ResponseWriter, req *http.Request) {
	srv.alog.Msg(req.RemoteAddr, code, req.Method, req.RequestURI, message)
	resp.Header().Add("content-type", "text/plain; charset=utf-8")
	resp.WriteHeader(code)
	resp.Write([]byte(message))
}

func (srv *Server) archiveHandler(
	fs filer.Filer,
	archive settings.ServedArchive,
	resp http.ResponseWriter,
	req *http.Request,
) {
	filePath, err := pathInArchive(archive, req.URL.Path)
	if err != nil {
		srv.writeError(404, "Not Found: "+err.Error(), resp, req)
		return
	}
	file, stat, err := fs.PreOpen(filePath)
	if err != nil {
		srv.writeError(404, "Not Found: "+err.Error(), resp, req)
		return
	}
	defer file.Close()

	var contentType string
	mime := fs.Mime(filePath)
	if mime.Type == "text" && srv.enc != "" {
		contentType = mime.Value + "; charset=" + srv.enc
	} else {
		contentType = mime.Value
	}

	resp.Header().Add("content-type", contentType)
	resp.Header().Add("content-length", fmt.Sprintf("%d", stat.Size()))
	resp.Header().Add("last-modified", HttpDate(stat.ModTime()))
	resp.Header().Add("x-path", filePath)
	resp.Header().Add("x-name", stat.Name())

	if !IsModifiedSince(req.Header.Get("if-modified-since"), stat.ModTime()) {
		srv.alog.Msg(req.RemoteAddr, 200, contentType, req.Method, req.RequestURI)
		resp.WriteHeader(304) // not modified
		return
	}

	srv.alog.Msg(req.RemoteAddr, 200, contentType, req.Method, req.RequestURI)
	resp.WriteHeader(200)

	if req.Method == http.MethodHead {
		return
	}

	// TODO needs benchmarking:
	//buf := bytes.NewBuffer(make([]byte, srv.bsize))
	buf := srv.bpool.New()
	defer srv.bpool.Return(buf)

	_, err = io.CopyBuffer(resp, file, buf.Bytes())
	if err != nil {
		srv.elog.Msg(err.Error())
		resp.Write([]byte(err.Error()))
		return
	}
}

func (srv *Server) ELog() logger.Logger {
	return srv.elog
}

func (srv *Server) ALog() logger.Logger {
	return srv.alog
}
