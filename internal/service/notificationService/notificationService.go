package notificationService

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"notification_service/config"
	"notification_service/internal/model"
	"notification_service/internal/notification/emailSender"
	"notification_service/internal/repository"
	"notification_service/internal/repository/notificationHistoryStatus"
	"notification_service/internal/utils"
	model2 "notification_service/pkg/notificationProducer/model"
	"regexp"
	"strings"
	"text/template"
)

type NotificationService struct {
	repo        repository.IRepository
	cfg         *config.Config
	emailSender emailSender.IEmailSender
}

func New(repo repository.IRepository, cfg *config.Config, emailSender emailSender.IEmailSender) *NotificationService {
	return &NotificationService{repo: repo, cfg: cfg, emailSender: emailSender}
}

func (s *NotificationService) Process(ctx context.Context, msg model2.NotificationMessage) error {
	// первично валидируем emailSender и templateName
	rqId := utils.GetRequestIdFromCtx(ctx)
	slog.Info("start process notification", slog.String("rqId", rqId), slog.Any("msg", msg))
	err := s.validateEmailAndTemplateName(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed on validate emailSender and template name: %w", err)
	}

	// идем в БД ищем темплейт и берем его обязательные поля
	notifTemplate, err := s.repo.GetTemplate(ctx, msg.TemplateName)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return errors.New("template not found")
		}
		return fmt.Errorf("failed on GetTemplate: %w", err)
	}

	// проверяем наличие обязательных полей
	err = s.validateRequiredParams(ctx, msg, notifTemplate)
	if err != nil {
		return fmt.Errorf("failed on validateRequiredParams: %w", err)
	}

	// рендерим шаблон
	renderedTmpl, err := s.renderTemplate(ctx, msg, notifTemplate)
	if err != nil {
		return fmt.Errorf("failed on renderTemplate: %w", err)
	}

	// отправляем email и записываем результат в БД
	err = s.emailSender.Send(msg.Email, msg.Subject, renderedTmpl)
	if err != nil {
		_ = s.saveIntoHistory(ctx, msg.Email, renderedTmpl, notificationHistoryStatus.FAILED, err.Error())
		return fmt.Errorf("failed on send email: %w", err)
	}

	_ = s.saveIntoHistory(ctx, msg.Email, renderedTmpl, notificationHistoryStatus.SUCCESS, "")
	return nil
}

func (s *NotificationService) saveIntoHistory(
	ctx context.Context,
	email, text, status, errMsg string,
) error {
	err := s.repo.InsertIntoHistory(ctx, email, text, status, errMsg)
	if err != nil {
		rqId := utils.GetRequestIdFromCtx(ctx)
		slog.Error(
			"got error from repo.InsertIntoHistory",
			slog.Any("err", err),
			slog.String("rqId", rqId),
		)
		return err
	}
	return nil
}

func (s *NotificationService) validateEmailAndTemplateName(ctx context.Context, msg model2.NotificationMessage) error {
	if msg.TemplateName == "" {
		return errors.New("template name is empty")
	}
	if msg.Email == "" {
		return errors.New("emailSender is empty")
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(msg.Email) {
		return errors.New("emailSender is invalid")
	}
	return nil
}

func (s *NotificationService) validateRequiredParams(
	ctx context.Context,
	msg model2.NotificationMessage,
	notifTemplate model.NotificationTemplate,
) error {
	for _, param := range notifTemplate.RequiredParameters {
		value, exists := msg.Parameters[param]
		if !exists {
			return fmt.Errorf("required parameter %s is missing", param)
		}
		if value == "" {
			return fmt.Errorf("required parameter %s is empty", param)
		}
	}
	return nil
}

func (s *NotificationService) renderTemplate(
	ctx context.Context,
	msg model2.NotificationMessage,
	notifTemplate model.NotificationTemplate,
) (string, error) {
	t, err := template.New("").Parse(notifTemplate.TemplateContent)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	err = t.Execute(&result, msg.Parameters)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}
