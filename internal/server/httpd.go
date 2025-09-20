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
	srv     *http.Server
	elog    logger.Logger
	alog    logger.Logger
	alogOff bool
	ctx     context.Context
	enc     string // default text encoding
	bsize   int
	bpool   *pool.BufferPool
	routes  []filer.Filer
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
		srv:    srv,
		ctx:    ctx,
		elog:   elog,
		alog:   alog,
		routes: make([]filer.Filer, 0, len(cfg.Archives)+len(cfg.Directories)),
	}
	res.srv.SetKeepAlivesEnabled(true)

	res.ApplyConfig(cfg)
	return res
}

func (srv *Server) ApplyConfig(cfg *settings.Settings) {
	srv.enc = cfg.DefaultChareset
	srv.bsize = cfg.BufferSize
	srv.bpool = pool.NewBufferPool(cfg.BufferSize)
	srv.alogOff = cfg.AccessLogOff

	mux := http.NewServeMux()
	rootHandled := false
	// arcives
	for _, archive := range cfg.Archives {
		if archive.UrlPrefix == "/" {
			rootHandled = true
		}
		handle, err := srv.newPrefixHandlerF(archive)
		if err != nil {
			srv.elog.Msg("Cant start handling file:", archive.ArchivePath, "[", err.Error(), "]")
			continue
		}
		mux.HandleFunc(archive.UrlPrefix, handle)
	}
	// arcive dirs
	for _, archiveDir := range cfg.Directories {
		if archiveDir.UrlPrefix == "/" {
			rootHandled = true
		}
		handle, err := srv.newPrefixHandlerD(archiveDir)
		if err != nil {
			srv.elog.Msg("Cant start handling dir:", archiveDir.DirectoryPath, "[", err.Error(), "]")
			continue
		}
		mux.HandleFunc(archiveDir.UrlPrefix, handle)
	}
	//
	if !rootHandled {
		mux.HandleFunc("/", srv.nullHandler)
	}
	srv.srv.Handler = mux

	timeout := time.Second * time.Duration(cfg.ClientTimeout)
	srv.srv.ReadTimeout = timeout
	srv.srv.WriteTimeout = timeout
}

func (srv *Server) newPrefixHandlerF(
	archive settings.ServedArchive,
) (func(resp http.ResponseWriter, req *http.Request), error) {
	fs, err := filer.NewFilerFromRoute(archive)
	if err != nil {
		return nil, err
	}
	srv.routes = append(srv.routes, fs)
	return func(resp http.ResponseWriter, req *http.Request) {
		srv.archiveHandler(fs, archive.UrlPrefix, resp, req)
	}, nil
}

func (srv *Server) newPrefixHandlerD(
	archive settings.ServedArchiveDir,
) (func(resp http.ResponseWriter, req *http.Request), error) {
	fs, err := filer.NewFilerDirFromRoute(archive)
	if err != nil {
		return nil, err
	}
	srv.routes = append(srv.routes, fs)
	return func(resp http.ResponseWriter, req *http.Request) {
		srv.archiveHandler(fs, archive.UrlPrefix, resp, req)
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

func (srv *Server) Release() {
	srv.elog.Msg("Releasing data files...")
	for _, fs := range srv.routes {
		fs.Release()
	}
}

func (srv *Server) nullHandler(resp http.ResponseWriter, req *http.Request) {
	srv.writeError(404, "Not Found", resp, req)
}

func (srv *Server) writeError(code int, message string, resp http.ResponseWriter, req *http.Request) {
	if !srv.alogOff {
		srv.alog.Msg(logFormatDefault(code, message, req)...)
	}
	if srv.alogOff && code != 404 {
		// log errors if it is not 404
		srv.elog.Msg(logFormatDefault(code, message, req)...)
	}
	resp.Header().Add("content-type", "text/plain; charset=utf-8")
	resp.WriteHeader(code)
	resp.Write([]byte(message))
}

func (srv *Server) archiveHandler(
	fs filer.Filer,
	urlPrefix string,
	resp http.ResponseWriter,
	req *http.Request,
) {
	filePath, err := pathUnderRoute(urlPrefix, req.URL.Path)
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
		if !srv.alogOff {
			srv.alog.Msg(logFormatDefault(304, "-", req)...)
		}
		resp.WriteHeader(304) // not modified
		return
	}

	if !srv.alogOff {
		srv.alog.Msg(logFormatDefault(200, "-", req)...)
	}
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
