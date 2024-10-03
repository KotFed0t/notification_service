package model

import "github.com/lib/pq"

type NotificationTemplate struct {
	Id                 int64          `db:"id"`
	TemplateName       string         `db:"template_name"`
	TemplateContent    string         `db:"template_content"`
	RequiredParameters pq.StringArray `db:"required_parameters"`
}
