package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/goexpert/lab-observabilidade-service-a/internal/infra/server"
	lab "github.com/goexpert/labobservabilidade"
)

func main() {

	slog.SetLogLoggerLevel(slog.LevelDebug)

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGINT)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer stop()

	otelShutdown, err := lab.InitProvider(ctx, "service-a", "otel-collector:4317")
	if err != nil {
		slog.Error("InitProvider", "error", err.Error())
		os.Exit(5)
	}
	defer func() {
		if err := otelShutdown(ctx); err != nil {
			slog.Error("otelShutdown", "error", err.Error())
			os.Exit(1)
		}
	}()

	webServer := lab.NewServer(os.Getenv("LO_PORT"))
	webServer.AddHandler("POST /cep", server.PostHandler)
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- webServer.Run()
	}()

	select {
	case <-srvErr:
		slog.Info("Serviço finalizando via <CTRL>+C...")
		return
	case <-ctx.Done():
		slog.Info("Serviço finalizando via interrupção no sistema.")
		stop()
	}
}
