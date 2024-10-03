package serviceInterface

import (
	"context"
	"notification_service/pkg/notificationProducer/model"
)

type INotificationService interface {
	Process(ctx context.Context, msg model.NotificationMessage) error
}
