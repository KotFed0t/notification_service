package emailSender

type IEmailSender interface {
	Send(to, subject, text string) error
}
