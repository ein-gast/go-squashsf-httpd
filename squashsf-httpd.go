package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/ein-gast/go-squashsf-httpd/internal/logger"
	"github.com/ein-gast/go-squashsf-httpd/internal/server"
	"github.com/ein-gast/go-squashsf-httpd/internal/settings"
)

func main() {
	config := settingsFromFlags()
	ctx := context.Background()
	log := logger.NewLogger()

	srv := server.NewServer(ctx, log, config)
	srv.Serve(log)
	log.Msg("App terminated")
}

func settingsFromFlags() *settings.Settings {
	config := settings.NewSettings()
	bindAddr := flag.String("host", config.BindAddr, "Bind this address")
	bindPort := flag.Int("port", config.BindPort, "Listen this port")
	prefix := flag.String("prefix", "/", "URL prefix")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Println("Set archive path")
		os.Exit(0)
	}
	config.BindAddr = *bindAddr
	config.BindPort = *bindPort
	config.Archives = append(
		config.Archives,
		settings.ServedArchive{
			UrlPrefix:   *prefix,
			ArchivePath: flag.Arg(0),
		},
	)
	return config
}
