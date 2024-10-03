package repository

import (
	"context"
	"notification_service/internal/model"
)

type IRepository interface {
	GetTemplate(ctx context.Context, templateName string) (model.NotificationTemplate, error)
	InsertIntoHistory(ctx context.Context, email, text, status, errorMessage string) error
}
