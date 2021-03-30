package paypal

import "testing"

func TestBalance(t *testing.T) {
	cli, err := NewClient("AVS9b9Ee9culv6VHec1r-4nb66MazgmGjZLQWLs0J4EuQV5Sco5uc7MDd8vW9YTB4BRnXu64yzLznmP5", "EGfVmGMuKcvwLsIgLgla42DHLXtPpiZR2SYfWvkRYx1rKmxYifarZYMMAbVhBfmCg_aRYRXsVCDQedgE", APIBaseSandBox)
	if err != nil {
		t.Error(err)
		return
	}

	out, err := cli.GetBalanceAccounts()
	if err != nil {
		t.Error(err)
		return
	}

	if out == nil {
		t.Error()
		return
	}
}
