package notificationProducer

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"notification_service/pkg/notificationProducer/model"
)

type NotificationProducer struct {
	writer *kafka.Writer
}

func New(kafkaUrl []string, topic string) *NotificationProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaUrl...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &NotificationProducer{writer: writer}
}

func (p *NotificationProducer) Send(ctx context.Context, key string, msg model.NotificationMessage) error {
	kafkaValue, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(key),
			Value: kafkaValue,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *NotificationProducer) Close() error {
	err := p.writer.Close()
	if err != nil {
		return err
	}

	return nil
}
