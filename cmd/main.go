package main

import (
	"github.com/KotFed0t/notification_service/config"
	"github.com/KotFed0t/notification_service/data/db/postgres"
	"github.com/KotFed0t/notification_service/data/queue/kafka/notificationConsumer"
	"github.com/KotFed0t/notification_service/internal/notification/emailSender"
	"github.com/KotFed0t/notification_service/internal/repository"
	"github.com/KotFed0t/notification_service/internal/service/notificationService"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	var logLevel slog.Level

	switch cfg.LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warning":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(log)

	slog.Debug("config", slog.Any("cfg", cfg))

	postgresDb := postgres.MustInitPostgres(cfg)
	repo := repository.NewPostgresRepo(postgresDb)

	mailSender := emailSender.New(cfg)

	notificationSrv := notificationService.New(repo, cfg, mailSender)

	notifConsumer := notificationConsumer.New(cfg, notificationSrv)
	go notifConsumer.Consume()

	// Waiting interruption signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	s := <-interrupt
	slog.Info("got interruption signal: " + s.String())

	notifConsumer.Close()

	err := postgresDb.Close()
	if err != nil {
		slog.Error("got error on postgresDb.Close()", slog.Any("err", err))
	}
}
