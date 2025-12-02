package channels

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/infra/mailer"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/templates"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/utils"
	"github.com/ritchieridanko/pasarly/backend/shared/ce"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/gomail.v2"
)

const emailErrTracer string = "channel.email"

type EmailChannel interface {
	SendWelcome(ctx context.Context, email, token string) (err error)
}

type emailChannel struct {
	baseURL  string
	sender   string
	mailer   *mailer.Mailer
	template *template.Template
}

func NewEmailChannel(m *mailer.Mailer, baseURL, sender string) (EmailChannel, error) {
	t, err := template.ParseFS(templates.FS, "*.html.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize email channel: %w", err)
	}

	return &emailChannel{baseURL: baseURL, sender: sender, mailer: m, template: t}, nil
}

func (c *emailChannel) SendWelcome(ctx context.Context, email, token string) error {
	_, span := otel.Tracer(emailErrTracer).Start(ctx, "SendWelcome")
	defer span.End()

	url, err := utils.URLWithToken(c.baseURL, "/auth/verify-account/confirm", token)
	if err != nil {
		e := fmt.Errorf("failed to send email: %w", err)
		utils.TraceErr(span, e, ce.MsgInternalServer)
		return e
	}

	data := struct {
		Email string
		URL   string
		Year  int
	}{
		Email: email,
		URL:   url,
		Year:  time.Now().UTC().Year(),
	}

	body, err := c.buildTemplate(span, "welcome.html.tmpl", data)
	if err != nil {
		return err
	}

	m := c.buildMessage([]string{email}, "Welcome to Pasarly!", body.String())
	return c.sendEmail(span, m)
}

func (c *emailChannel) buildTemplate(s trace.Span, template string, data any) (bytes.Buffer, error) {
	var b bytes.Buffer
	if err := c.template.ExecuteTemplate(&b, template, data); err != nil {
		e := fmt.Errorf("failed to send email: %w", err)
		utils.TraceErr(s, e, ce.MsgInternalServer)
		return bytes.Buffer{}, e
	}

	return b, nil
}

func (c *emailChannel) buildMessage(recipients []string, subject, body string) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", c.sender)
	m.SetHeader("To", recipients...)
	m.SetHeader("Subject", utils.MIMEBase64(subject))
	m.SetBody("text/plain", "Please view this email in an HTML-compatible client!")
	m.AddAlternative("text/html", body)
	return m
}

func (c *emailChannel) sendEmail(s trace.Span, m *gomail.Message) error {
	if err := c.mailer.Send(m); err != nil {
		e := fmt.Errorf("failed to send email: %w", err)
		utils.TraceErr(s, e, ce.MsgInternalServer)
		return e
	}
	return nil
}
