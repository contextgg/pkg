package types

import "testing"

type testIt struct{}
type testIt2 struct{}

func Test_It(t *testing.T) {
	t1 := &testIt{}
	tf := func() interface{} {
		return &testIt2{}
	}

	reg := NewRegistry()
	reg.Add(EntryFromType(t1, true))
	reg.Add(EntryFromFactory(tf, true))

	entry, ok := reg.GetByName("testit")
	if !ok {
		t.Error("Whoops")
		return
	}
	if entry == nil {
		t.Error("Whoops")
	}
}
