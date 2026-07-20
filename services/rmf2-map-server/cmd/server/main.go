package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rmf2map/server/internal/api"
	"rmf2map/server/internal/config"
	"rmf2map/server/internal/middleware"
	"rmf2map/server/internal/service"
	"rmf2map/server/internal/vda5050"
)

func main() {
	cfg := config.Config{DefaultOrg: "ros-industrial"}

	flag.StringVar(&cfg.Addr, "addr", "0.0.0.0:8008", "listen address (host:port)")
	flag.StringVar(&cfg.DataDir, "data-dir", "./data", "path to data directory")
	flag.StringVar(&cfg.DefaultOrg, "default-org", cfg.DefaultOrg, "default organisation")
	flag.StringVar(&cfg.VDA5050Base, "vda5050-base", "", "VDA5050 master base URL (e.g. http://localhost:8000); empty disables live positions")
	flag.Parse()

	keys, err := service.LoadAPIKeys(cfg.DataDir)
	if err != nil {
		log.Fatalf("load api keys: %v", err)
	}

	// pollerCtx is cancelled on shutdown so the background goroutine exits cleanly.
	pollerCtx, cancelPoller := context.WithCancel(context.Background())
	defer cancelPoller()

	var poller *vda5050.Poller
	if cfg.VDA5050Base != "" {
		poller = vda5050.NewPoller(cfg.VDA5050Base)
		poller.Start(pollerCtx)
		log.Printf("vda5050 live positions enabled, polling %s/state every 200ms", cfg.VDA5050Base)
	} else {
		log.Printf("vda5050 not configured — positions served from positions.json")
	}

	mux := http.NewServeMux()
	api.RegisterRoutes(mux, cfg, poller)

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           middleware.New(keys, mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("listening on %s, data dir: %s, default org: %s", cfg.Addr, cfg.DataDir, cfg.DefaultOrg)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down...")
	cancelPoller()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
	log.Println("server stopped")
}
