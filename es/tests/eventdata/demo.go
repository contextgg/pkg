package eventdata

import "github.com/contextgg/pkg/types"

type DemoCreated struct {
	Name string
}
type DemoDescriptionAdded struct {
	Description string
}

func init() {
	types.SetTypeData(&DemoCreated{}, true)
	types.SetTypeData(&DemoDescriptionAdded{}, true)
}
