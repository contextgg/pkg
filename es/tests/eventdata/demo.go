package eventdata

import "github.com/contextgg/pkg/types"

type DemoCreated struct {
	Name string
}
type DemoDescriptionAdded struct {
	Description string
}

func init() {
	types.Add(types.EntryFromType(&DemoCreated{}, true))
	types.Add(types.EntryFromType(&DemoDescriptionAdded{}, true))
}
