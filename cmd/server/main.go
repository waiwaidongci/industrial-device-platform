package main

import (
	"context"
	"errors"
	"github.com/example/industrial-device-platform/internal"
	deviceadapter "github.com/example/industrial-device-platform/internal/devicecatalog/adapter"
	"github.com/example/industrial-device-platform/internal/devicecatalog/application"
	"github.com/example/industrial-device-platform/internal/devicecatalog/infrastructure"
	commandadapter "github.com/example/industrial-device-platform/internal/remotecommand/adapter"
	commandapp "github.com/example/industrial-device-platform/internal/remotecommand/application"
	commandinfra "github.com/example/industrial-device-platform/internal/remotecommand/infrastructure"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := internal.LoadConfig("configs/config.yaml")
	log := internal.NewLogger()
	metrics := &internal.Metrics{}
	devices := infrastructure.NewMemoryRepository()
	deviceSvc := application.NewService(devices)
	deviceHandler := deviceadapter.NewHandler(deviceSvc, metrics)
	commands := commandinfra.NewMemoryRepository()
	commandSvc := commandapp.NewService(commands)
	commandHandler := commandadapter.NewHandler(commandSvc, metrics)
	mux := http.NewServeMux()
	mux.Handle("/v1/devices", deviceHandler)
	mux.Handle("/v1/devices/", deviceHandler)
	mux.Handle("/v1/commands", commandHandler)
	mux.Handle("/v1/commands/", commandHandler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { internal.JSON(w, 200, map[string]string{"status": "ok"}) })
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		internal.JSON(w, 200, map[string]string{"status": "ready"})
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) { metrics.Render(w) })
	h := internal.Recover(log, metrics, internal.RequestID(internal.AccessLog(log, metrics, internal.CORS(internal.Timeout(mux)))))
	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: h, ReadHeaderTimeout: 5 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		log.Info("server started", "addr", cfg.HTTPAddr)
		if e := srv.ListenAndServe(); e != nil && !errors.Is(e, http.ErrServerClosed) {
			log.Error("server stopped", "error", e)
		}
	}()
	<-ctx.Done()
	shut, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	_ = srv.Shutdown(shut)
	log.Info("server shutdown")
}
