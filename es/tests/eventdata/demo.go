package eventdata

import "github.com/contextgg/pkg/types"

type DemoCreated struct {
	Name string
}
type DemoDescriptionAdded struct {
	Description string
}

func init() {
	types.Add(types.RegisterFromType(&DemoCreated{}), types.IsInternalType)
	types.Add(types.RegisterFromType(&DemoDescriptionAdded{}), types.IsInternalType)
}
