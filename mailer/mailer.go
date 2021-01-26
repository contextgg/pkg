package mailer

import (
	"errors"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Mailer interface {
	SendTemplate(toAddress, toName, templateId string, data map[string]interface{}) error
}

type mailer struct {
	cli  *sendgrid.Client
	from *mail.Email
}

func (s *mailer) SendTemplate(toAddress, toName, templateId string, data map[string]interface{}) error {
	m := mail.NewV3Mail()
	m.SetFrom(s.from)
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

func NewSendGrid(key string) Mailer {
	cli := sendgrid.NewSendClient(key)
	from := mail.NewEmail("Inflow", "no-reply@inflow.pro")
	return &mailer{cli, from}
}
