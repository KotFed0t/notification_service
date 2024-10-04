package repository

import (
	"context"
	"github.com/KotFed0t/notification_service/internal/model"
)

type IRepository interface {
	GetTemplate(ctx context.Context, templateName string) (model.NotificationTemplate, error)
	InsertIntoHistory(ctx context.Context, email, templateName, text, status, errorMessage string) error
}
