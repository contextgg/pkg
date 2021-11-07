package eventdata

import "github.com/contextgg/pkg/types"

type DemoCreated struct {
	Name string
}

func init() {
	types.SetTypeData(&DemoCreated{}, true)
}
