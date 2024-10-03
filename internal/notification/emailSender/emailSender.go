package emailSender

import (
	"gopkg.in/gomail.v2"
	"notification_service/config"
)

type EmailSender struct {
	dialer *gomail.Dialer
	cfg    *config.Config
}

func New(cfg *config.Config) *EmailSender {
	dialer := gomail.NewDialer(cfg.Mail.Host, cfg.Mail.Port, cfg.Mail.Address, cfg.Mail.Password)
	return &EmailSender{dialer: dialer, cfg: cfg}
}

func (es *EmailSender) Send(to, subject, text string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", es.cfg.Mail.Address)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", text)

	if err := es.dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
