package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ein-gast/go-squashsf-httpd/internal/filer"
	"github.com/ein-gast/go-squashsf-httpd/internal/logger"
	"github.com/ein-gast/go-squashsf-httpd/internal/server"
	"github.com/ein-gast/go-squashsf-httpd/internal/settings"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	config := settingsFromFlags()
	elog := logger.NewLogger()
	alog := logger.NewLogger()
	reopenLogs(elog, alog, config)

	elog.Msg("Starting with configuation")
	settings.PrintSetting(*config, elog)

	elog.Msg("Adding MIME types...")
	filer.AddMimeTypes()

	elog.Msg("PID:", os.Getpid())

	srv := server.NewServer(ctx, elog, alog, config)
	elog.Msg("Installing signal hook...")
	go hookSignal(ctx, cancel, elog, srv, config)

	srv.Serve()
	elog.Msg("App terminated")
}

func settingsFromFlags() *settings.Settings {
	config := settings.NewSettings()
	var err error
	yamlPath := flag.String("config", "", "Config file path")
	bindAddr := flag.String("host", config.BindAddr, "Bind this address")
	bindPort := flag.Int("port", config.BindPort, "Listen this port")
	prefix := flag.String("prefix", "/", "URL prefix")
	squash := flag.String("squash", "", "SquashFS file path")
	charset := flag.String("charset", config.DefaultChareset, "Default charset for text")
	flag.Parse()

	if *yamlPath != "" {
		config, err = settings.Load(*yamlPath)
		if err != nil {
			fmt.Println("Config reading error:", err.Error())
			os.Exit(0)
		}
	} else {
		config.BindAddr = *bindAddr
		config.BindPort = *bindPort
		config.DefaultChareset = *charset
	}

	if len(config.Archives) == 0 && *squash == "" {
		fmt.Println("At least one SquashFS file path must be provided in CLI or config")
		os.Exit(0)
	}
	if *squash != "" {
		config.Archives = append(
			config.Archives,
			settings.ServedArchive{
				UrlPrefix:   *prefix,
				ArchivePath: *squash,
			},
		)
	}
	return config
}

func reopenLogs(e logger.Logger, a logger.Logger, s *settings.Settings) {
	e.Msg("Reopening log:", s.ErrorLog)
	err := e.OpenFile(s.ErrorLog)
	if err != nil {
		e.Msg(err.Error())
	} else {
		e.Msg("Log opened: ", s.ErrorLog)
	}
	e.Msg("Reopening log:", s.AccessLog)
	err = a.OpenFile(s.AccessLog)
	if err != nil {
		e.Msg(err.Error())
	} else {
		e.Msg("Log opened: ", s.AccessLog)
	}
}

func hookSignal(
	ctx context.Context,
	cancel context.CancelFunc,
	log logger.Logger,
	srv *server.Server,
	s *settings.Settings,
) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1)

	for {
		select {
		case sig := <-c:
			log.Msg("Got signal:", sig)
			switch sig {
			case syscall.SIGUSR1:
				log.Msg("Reloading by signal...")
				reopenLogs(srv.ELog(), srv.ALog(), s)
				srv.Release()
			default:
				log.Msg("Terminaging by signal...")
				srv.Shutdown()
				cancel()
				return
			}
		case <-ctx.Done():
			log.Msg("Canceled")
			return
		}
	}
}
