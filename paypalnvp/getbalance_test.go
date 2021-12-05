package paypalnvp

import "testing"

func TestIt(t *testing.T) {
	cli := NewClient("sb-dmlkb59712_api1.business.example.com", "XZ88CU2T7HU3ZMQP", "AWtQMwn7D2rcuzAnqPZZmP7IndbwA4nLRVMCEAWyO0WUqrTK1SFxWzQ.", true)
	amounts, err := cli.GetBalance(true)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(amounts)
}
