package notificationProducer

import (
	"context"
	"github.com/KotFed0t/notification_service/pkg/notificationProducer/model"
)

type INotificationProducer interface {
	Send(ctx context.Context, key string, msg model.NotificationMessage) error
}
