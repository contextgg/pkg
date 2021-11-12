package types

type TypeOption func(o *Entry)

var IsInternalType = TypeOption(func(e *Entry) {
	e.InternalType = true
})

func RegisterFromType(obj interface{}) TypeOption {
	// get it!
	t := getElemType(obj)
	name := getShortName(t)

	return func(e *Entry) {
		e.Name = name
		e.Fullname = t.String()
		e.Type = t
	}
}

func RegisterFromFactory(factory TypeFunc) TypeOption {
	opt := RegisterFromType(factory())
	return func(e *Entry) {
		opt(e)
		e.Factory = factory
	}
}
