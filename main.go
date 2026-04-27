// Command numbers serves the current "number of the week" derived from book.txt.
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/icco/gutil/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/zap"
)

// main wires dependencies and blocks until SIGINT/SIGTERM.
func main() {
	log, err := logging.NewLogger("numbers")
	if err != nil {
		fallback, ferr := zap.NewProduction()
		if ferr != nil {
			fmt.Fprintf(os.Stderr, "logger init: %v / %v\n", err, ferr)
			os.Exit(1)
		}
		fallback.Warn("falling back to zap.NewProduction logger", zap.Error(err))
		log = fallback.Sugar()
	}
	defer func() {
		if err := log.Sync(); err != nil {
			log.Debugw("logger sync", zap.Error(err))
		}
	}()

	registry := prometheus.NewRegistry()
	exporter, err := otelprom.New(otelprom.WithRegisterer(registry))
	if err != nil {
		log.Errorw("otel prometheus exporter", zap.Error(err))
		os.Exit(1)
	}
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	otel.SetMeterProvider(mp)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := mp.Shutdown(shutdownCtx); err != nil {
			log.Warnw("meter provider shutdown", zap.Error(err))
		}
	}()

	port := "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	srv := &http.Server{
		Addr: ":" + port,
		Handler: Router(Options{
			Logger:         log,
			MetricsHandler: promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
		}),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx = logging.NewContext(ctx, log)

	go func() {
		log.Infow("http server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorw("http server", zap.Error(err))
			stop()
		}
	}()

	<-ctx.Done()
	log.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Errorw("http shutdown", zap.Error(err))
	}
}
