package mailer

import (
	"os"
	"testing"
)

func TestEmail(t *testing.T) {
	key := os.Getenv("SENDGRID_KEY")

	m := NewSendGrid(key, "no-reply@inflow.services")

	if err := m.SendTemplate("Me", "nathan@context.gg", "sHadey", "d-5f41f02cc5f9412aac6d9f775a1c90bc", nil); err != nil {
		t.Error(err)
		return
	}
	if err := m.SendTemplate("Me", "chris@context.gg", "doofy", "d-5f41f02cc5f9412aac6d9f775a1c90bc", nil); err != nil {
		t.Error(err)
		return
	}
}
