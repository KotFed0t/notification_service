package notificationProducer

import (
	"context"
	"notification_service/pkg/notificationProducer/model"
)

type INotificationProducer interface {
	Send(ctx context.Context, key string, msg model.NotificationMessage) error
}
