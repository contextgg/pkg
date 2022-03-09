package es

import "github.com/contextgg/pkg/types"

var publicEvents = map[string]bool{}

func AddPublicEvent(obj interface{}) {
	name := types.GetTypeName(obj)
	publicEvents[name] = true
}

func IsPublicEvent(name string) bool {
	_, ok := publicEvents[name]
	return ok
}
