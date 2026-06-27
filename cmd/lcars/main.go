// Command lcars runs the LCARS engine HTTP server.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adejong5/lcars-engine/internal/config"
	"github.com/adejong5/lcars-engine/internal/ha"
	"github.com/adejong5/lcars-engine/internal/render"
	"github.com/adejong5/lcars-engine/internal/server"
)

func main() {
	cfg := config.Load()

	level := slog.LevelInfo
	if cfg.Dev {
		level = slog.LevelDebug
	}
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	// Data source: mock fabricated data, or a live HA WebSocket connection.
	// Both satisfy ha.Source, so the HTTP layer is identical either way.
	var src ha.Source
	var startSource func(context.Context)
	if cfg.Mock {
		mock := ha.NewMock()
		src = mock
		startSource = func(ctx context.Context) { mock.Run(ctx, 2*time.Second) }
	} else {
		if cfg.HAToken == "" || (!cfg.Addon && cfg.HAHost == "") {
			log.Error("MOCK=false requires HA_TOKEN (and HA_HOST unless ADDON=true)")
			os.Exit(1)
		}
		wsURL := ha.BuildWSURL(cfg.HAHost, cfg.HASSL)
		if cfg.Addon {
			wsURL = ha.SupervisorWSURL // reach HA through the Supervisor proxy
		}
		live := ha.NewLive(wsURL, cfg.HAToken, log)
		src = live
		startSource = func(ctx context.Context) { live.Run(ctx) }
	}

	rnd, err := render.New(cfg.Dev)
	if err != nil {
		log.Error("parse templates", "err", err)
		os.Exit(1)
	}

	srv := server.New(cfg, log, src, rnd)
	httpServer := &http.Server{
		Addr:              cfg.Addr,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Serve until SIGINT/SIGTERM, then shut down gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start the data source: the mock drift loop, or the live HA connection.
	go startSource(ctx)

	go func() {
		log.Info("listening", "addr", cfg.Addr, "mock", cfg.Mock, "dev", cfg.Dev)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server failed", "err", err)
			stop()
		}
	}()

	<-ctx.Done()
	log.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", "err", err)
		os.Exit(1)
	}
}
