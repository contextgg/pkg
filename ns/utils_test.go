package ns

import (
	"net/http"
	"testing"
)

func TestIt(t *testing.T) {
	data := []struct {
		in  string
		out string
	}{
		{in: "localhost:3000", out: "localhost"},
		{in: "localhost", out: "localhost"},
		{in: "inflow.pro", out: "inflowpro"},
		{in: "contextgg.inflow.pro", out: "contextgg"},
		{in: "abc.inflow.pro", out: "abc"},
	}

	for _, d := range data {
		t.Run(d.in, func(t *testing.T) {
			r := &http.Request{
				Host: d.in,
			}

			b := Slug(r, ".inflow.pro")
			if b != d.out {
				t.Errorf("Slug mismatch %s != %s", b, d.out)
			}
		})
	}
}
