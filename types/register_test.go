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
	reg.Add(RegisterFromType(t1))
	reg.Add(RegisterFromFactory(tf))

	entry, ok := reg.GetByName("testit")
	if !ok {
		t.Error("Whoops")
		return
	}

	if entry == nil {
		t.Error("Whoops")
	}
}
