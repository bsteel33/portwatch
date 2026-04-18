package daemon

import (
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

// Daemon periodically scans ports and alerts on changes.
type Daemon struct {
	cfg     *config.Config
	scanner *scanner.Scanner
	alerter *alert.Alerter
}

// New creates a new Daemon with the given config.
func New(cfg *config.Config) *Daemon {
	return &Daemon{
		cfg:     cfg,
		scanner: scanner.New(cfg),
		alerter: alert.New(cfg),
	}
}

// Run starts the daemon loop, blocking until stop is closed.
func (d *Daemon) Run(stop <-chan struct{}) error {
	if err := d.tick(); err != nil {
		log.Printf("[portwatch] initial scan error: %v", err)
	}

	ticker := time.NewTicker(time.Duration(d.cfg.Interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := d.tick(); err != nil {
				log.Printf("[portwatch] scan error: %v", err)
			}
		case <-stop:
			log.Println("[portwatch] daemon stopped")
			return nil
		}
	}
}

func (d *Daemon) tick() error {
	ports, err := d.scanner.OpenPorts()
	if err != nil {
		return err
	}

	current := snapshot.New(ports)

	prev, err := snapshot.Load(d.cfg.SnapshotPath)
	if err != nil {
		log.Printf("[portwatch] no previous snapshot, saving baseline")
		return current.Save(d.cfg.SnapshotPath)
	}

	diff := snapshot.Compare(prev, current)
	if err := d.alerter.Notify(diff); err != nil {
		log.Printf("[portwatch] alert error: %v", err)
	}

	return current.Save(d.cfg.SnapshotPath)
}
