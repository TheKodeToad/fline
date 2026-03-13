package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	fine "github.com/TheKodeToad/fine/internal"
	"github.com/TheKodeToad/fine/internal/api"
	"github.com/TheKodeToad/fine/internal/config"
	"github.com/TheKodeToad/fine/internal/gateway"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", slog.Any("err", err))
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: conf.LogLevel,
	}))
	slog.SetDefault(logger)

	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)

	router.Mount("/api", api.Routes(&conf))

	var gateway gateway.Gateway
	router.Get(fine.GatewayPath, func(w http.ResponseWriter, r *http.Request) {
		gateway.ServeHTTP(&conf, w, r)
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("no route for " + r.URL.Path)
		http.NotFound(w, r)
	})

	server := http.Server{
		Addr: conf.ListenAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Debug("handling " + r.URL.Path)
			router.ServeHTTP(w, r)
		}),
	}

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		slog.Error("net.Listen failed", slog.Any("err", err))
		os.Exit(1)
	}

	go func() {
		err = server.Serve(listener)
		if err != http.ErrServerClosed {
			slog.Error("error serving connections", slog.Any("err", err))
		}
	}()

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGTERM, syscall.SIGINT)

	<-exitSignal

	slog.Info("goodbye")

	slog.Info("shutting down HTTP server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	server.Shutdown(ctx)
	cancel()

	slog.Info("shutting down gateway")
	gateway.Shutdown()
}
