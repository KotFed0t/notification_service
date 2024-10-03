package notificationConsumer

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/KotFed0t/notification_service/config"
	"github.com/KotFed0t/notification_service/internal/service/serviceInterface"
	"github.com/KotFed0t/notification_service/internal/utils"
	"github.com/KotFed0t/notification_service/pkg/notificationProducer/model"
	"github.com/segmentio/kafka-go"
	"io"
	"log/slog"
)

type NotificationConsumer struct {
	reader          *kafka.Reader
	cfg             *config.Config
	notificationSrv serviceInterface.INotificationService
}

func New(
	cfg *config.Config,
	notificationSrv serviceInterface.INotificationService,
) *NotificationConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.KafkaNotification.ConsumerUrl,
		Topic:   cfg.KafkaNotification.Topic,
		GroupID: cfg.KafkaNotification.ConsumerGroup,
	})

	return &NotificationConsumer{
		reader:          reader,
		cfg:             cfg,
		notificationSrv: notificationSrv,
	}
}

func (c *NotificationConsumer) Consume() {
	for {
		m, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			if errors.Is(err, io.EOF) {
				slog.Info("reader has been closed, stop consuming.")
				break
			}
			slog.Error("error while ReadMessage from kafka notification", slog.Any("error", err))
		}
		go c.handleMessage(m)
	}
}

func (c *NotificationConsumer) handleMessage(m kafka.Message) {
	var msg model.NotificationMessage
	err := json.Unmarshal(m.Value, &msg)
	if err != nil {
		slog.Error("error while unmarshal json NotificationConsumer.handleMessage", slog.Any("error", err))
	}

	ctx := utils.CreateCtxWithRqID(context.Background())
	rqId := utils.GetRequestIdFromCtx(ctx)
	slog.Info(
		"got kafka notification message",
		slog.Any("msg", msg),
		slog.String("rqId", rqId),
	)

	err = c.notificationSrv.Process(ctx, msg)
	if err != nil {
		slog.Error(
			"got error from notificationSrv.Process",
			slog.Any("error", err),
			slog.String("rqId", rqId),
		)
	} else {
		slog.Info(
			"kafka message successfully processed",
			slog.String("rqId", rqId),
		)
	}
}

func (c *NotificationConsumer) Close() {
	err := c.reader.Close()
	if err != nil {
		slog.Error("error while closing kafka reader", slog.Any("error", err))
	}
}
