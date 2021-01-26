package gen

import "testing"

func TestIt(t *testing.T) {
	g := RandString(10)
	if len(g) != 10 {
		t.Error("Wrong len")
		return
	}
}
