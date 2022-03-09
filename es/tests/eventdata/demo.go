package eventdata

import "github.com/contextgg/pkg/types"

type DemoCreated struct {
	Name string
}
type DemoDescriptionAdded struct {
	Description string
}

func init() {
	entries := []*types.Entry{
		types.NewEntryFromObject(&DemoCreated{}),
		types.NewEntryFromObject(&DemoDescriptionAdded{}),
	}

	types.Add(entries...)
}
