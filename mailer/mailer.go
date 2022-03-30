package mailer

import (
	"errors"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Mailer interface {
	SendTemplate(fromName string, toAddress string, toName string, templateId string, data map[string]interface{}) error
}

type mailer struct {
	cli       *sendgrid.Client
	fromEmail string
}

func (s *mailer) SendTemplate(fromName string, toAddress string, toName string, templateId string, data map[string]interface{}) error {
	from := mail.NewEmail(fromName, s.fromEmail)

	m := mail.NewV3Mail()
	m.SetFrom(from)
	m.SetTemplateID(templateId)

	p := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail(toName, toAddress),
	}
	p.AddTos(tos...)

	for key, value := range data {
		p.SetDynamicTemplateData(key, value)
	}
	m.AddPersonalizations(p)

	resp, err := s.cli.Send(m)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New("Wrong statuscode")
	}
	return nil
}

func NewSendGrid(apiKey string, fromEmail string) Mailer {
	cli := sendgrid.NewSendClient(apiKey)
	return &mailer{cli, fromEmail}
}
