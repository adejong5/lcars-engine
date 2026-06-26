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
		if cfg.HAHost == "" || cfg.HAToken == "" {
			log.Error("MOCK=false requires HA_HOST and HA_TOKEN")
			os.Exit(1)
		}
		live := ha.NewLive(cfg.HAHost, cfg.HAToken, cfg.HASSL, log)
		src = live
		startSource = func(ctx context.Context) { live.Run(ctx) }
	}

	srv := server.New(cfg, log, src)
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
