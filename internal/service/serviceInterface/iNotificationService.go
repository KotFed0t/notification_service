package serviceInterface

import (
	"context"
	"github.com/KotFed0t/notification_service/internal/model"
)

type INotificationService interface {
	Process(ctx context.Context, msg model.NotificationMessage) error
}
