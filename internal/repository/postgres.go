package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/KotFed0t/notification_service/internal/model"
	"github.com/KotFed0t/notification_service/internal/utils"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type PostgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db}
}

func (r *PostgresRepo) GetTemplate(ctx context.Context, templateName string) (
	model.NotificationTemplate,
	error,
) {
	var notificationTemplate model.NotificationTemplate
	query := `SELECT * FROM templates WHERE template_name = $1;`

	rqId := utils.GetRequestIdFromCtx(ctx)
	slog.Info("GetTemplate start", slog.String("rqId", rqId))

	err := r.db.GetContext(ctx, &notificationTemplate, query, templateName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Info(
				"GetTemplate not found",
				slog.String("rqId", rqId),
				slog.String("template_name", templateName),
			)
			return model.NotificationTemplate{}, ErrNoRows
		}

		slog.Error("GetTemplate error", slog.String("rqId", rqId), slog.Any("err", err))
		return model.NotificationTemplate{}, err
	}

	slog.Info("GetTemplate success", slog.String("rqId", rqId))
	return notificationTemplate, nil
}

func (r *PostgresRepo) InsertIntoHistory(
	ctx context.Context,
	email, templateName, text, status, errorMessage string,
) error {
	query := `INSERT INTO notification_history(email, template_name, text, status, error_message) VALUES ($1, $2, $3, $4, $5);`

	rqId := utils.GetRequestIdFromCtx(ctx)
	slog.Info("InsertIntoHistory start", slog.String("rqId", rqId))

	_, err := r.db.ExecContext(ctx, query, email, templateName, text, status, errorMessage)
	if err != nil {
		slog.Error("InsertIntoHistory error", slog.String("rqId", rqId))
		return err
	}

	slog.Info("InsertIntoHistory success", slog.String("rqId", rqId))
	return nil
}
