package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/ein-gast/go-squashsf-httpd/internal/daemon"
	"github.com/ein-gast/go-squashsf-httpd/internal/filer"
	"github.com/ein-gast/go-squashsf-httpd/internal/logger"
	"github.com/ein-gast/go-squashsf-httpd/internal/server"
	"github.com/ein-gast/go-squashsf-httpd/internal/settings"
)

var Version string = "from-sources"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	config := settingsFromFlags()
	elog := logger.NewLogger()
	alog := logger.NewLogger()
	reopenLogs(elog, alog, config)

	elog.Msg("Starting with configuation")
	settings.PrintSetting(*config, Version, elog)

	elog.Msg("Adding MIME types...")
	filer.AddMimeTypes()

	pid := daemon.Pid(os.Getpid())
	elog.Msg("PID:", pid)
	createPid(pid, elog, config)

	srv := server.NewServer(ctx, elog, alog, config)
	elog.Msg("Installing signal hook...")
	go hookSignal(ctx, cancel, elog, srv, config)

	srv.Serve()
	removePid(pid, elog, config)
	elog.Msg("App terminated")
}

func settingsFromFlags() *settings.Settings {
	config := settings.NewSettings()
	var err error
	version := flag.Bool("version", false, "Show version")
	yamlPath := flag.String("config", "", "Config file path")
	bindAddr := flag.String("host", config.BindAddr, "Bind this address")
	bindPort := flag.Int("port", config.BindPort, "Listen this port")
	prefix := flag.String("prefix", "/", "URL prefix")
	squash := flag.String("squash", "", "SquashFS file path")
	charset := flag.String("charset", config.DefaultChareset, "Default charset for text")
	flag.Parse()

	if *version {
		fmt.Println("Version:", Version)
		fmt.Println("Golang:", runtime.Version())
		os.Exit(0)
	}

	if *yamlPath != "" {
		config, err = settings.Load(*yamlPath)
		if err != nil {
			fmt.Println("Config reading error:", err.Error())
			os.Exit(1)
		}
	} else {
		config.BindAddr = *bindAddr
		config.BindPort = *bindPort
		config.DefaultChareset = *charset
	}

	if len(config.Archives) == 0 && *squash == "" {
		fmt.Println("At least one SquashFS file path must be provided in CLI or config")
		os.Exit(1)
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
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

	for {
		select {
		case sig := <-c:
			log.Msg("Got signal:", sig)
			switch sig {
			case syscall.SIGUSR1:
				log.Msg("Reopening logs by signal...")
				reopenLogs(srv.ELog(), srv.ALog(), s)
			case syscall.SIGUSR2:
				log.Msg("Reopening archives signal...")
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

func createPid(pid daemon.Pid, log logger.Logger, cfg *settings.Settings) {
	if cfg.PidFileOff {
		log.Msg("PID file is off - skipping creation")
		return
	}
	xpid, err := daemon.WritePidFileIfAbsent(pid, cfg, false)
	if err == nil {
		log.Msg("PID file created")
		return
	}
	if err == daemon.E_PID_EXIST {
		p, err := os.FindProcess(int(xpid))
		if err == nil && p.Signal(syscall.Signal(0)) == nil {
			log.Msg("The PID file of running process exists:", xpid)
			log.Msg("Refusing to start due to PID error")
			os.Exit(1)
		}
	}
	_, err = daemon.WritePidFileIfAbsent(pid, cfg, true)
	if err != nil {
		log.Msg(err.Error())
		log.Msg("Refusing to start due to PID error")
		os.Exit(1)
	}
	log.Msg("PID file owerwritten")
}

func removePid(pid daemon.Pid, log logger.Logger, cfg *settings.Settings) {
	if cfg.PidFileOff {
		log.Msg("PID file is off - skipping removing")
		return
	}
	xpid, err := daemon.RemovePidFile(pid, cfg, false)
	if err == nil {
		log.Msg("PID file removed")
		return
	}
	if err == daemon.E_PID_IS_NOT_MINE {
		p, err := os.FindProcess(int(xpid))
		if err == nil && p.Signal(syscall.Signal(0)) == nil {
			log.Msg("Note: PID file exists and points to another RUNNING process")
			log.Msg("PID file stays untouched")
		} else {
			log.Msg("Note: PID file exists and points to another DEAD process")
		}
	}
	_, err = daemon.RemovePidFile(pid, cfg, true)
	if err != nil {
		log.Msg(err.Error())
		return
	}
	log.Msg("PID file removed by force")
}
