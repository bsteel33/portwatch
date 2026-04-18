package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/daemon"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		log.Printf("using default config: %v", err)
		cfg = config.Default()
	}

	flags := config.ParseFlags()
	config.Apply(cfg, flags)

	log.Printf("[portwatch] starting — interval=%ds snapshot=%s",
		cfg.Interval, cfg.SnapshotPath)

	stop := make(chan struct{})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		close(stop)
	}()

	d := daemon.New(cfg)
	if err := d.Run(stop); err != nil {
		log.Fatalf("[portwatch] fatal: %v", err)
	}
}
