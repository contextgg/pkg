package types

import "testing"

type testIt struct{}
type testIt2 struct{}

func Test_It(t *testing.T) {
	t1 := &testIt{}

	reg := NewRegistry()
	reg.Upsert(t1)

	entry, ok := reg.GetByName("testit")
	if !ok {
		t.Error("Whoops")
		return
	}
	if entry == nil {
		t.Error("Whoops")
	}
}
