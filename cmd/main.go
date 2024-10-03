package main

import (
	"context"
	"log/slog"
	"notification_service/config"
	"notification_service/data/db/postgres"
	"notification_service/data/queue/kafka/notificationConsumer"
	"notification_service/internal/notification/emailSender"
	"notification_service/internal/repository"
	"notification_service/internal/service/notificationService"
	"notification_service/pkg/notificationProducer"
	"notification_service/pkg/notificationProducer/model"
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

	notifProducer := notificationProducer.New(cfg.KafkaNotification.ConsumerUrl, cfg.KafkaNotification.Topic)

	err := notifProducer.Send(context.Background(), "", model.NotificationMessage{
		Email:        "Slayvi555@gmail.com",
		Subject:      "test5",
		TemplateName: "test",
		Parameters: map[string]string{
			"Name": "Слава",
		},
	})
	if err != nil {
		slog.Error("got error from notifProducer.Send", slog.Any("err", err))
	}

	// Waiting interruption signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	s := <-interrupt
	slog.Info("got interruption signal: " + s.String())
	notifConsumer.Close()
	err = postgresDb.Close()
	if err != nil {
		slog.Error("got error on postgresDb.Close()", slog.Any("err", err))
	}
}
