package types

func EntryFromType(obj interface{}, internalType bool, names ...string) *Entry {
	// get it!
	t := getElemType(obj)
	n := getShortName(t)
	all := append([]string{n}, names...)

	return &Entry{
		Names:        all,
		Type:         t,
		InternalType: internalType,
		Factory:      typeFactory(t),
	}
}

func EntryFromFactory(factory TypeFunc, internalType bool, names ...string) *Entry {
	obj := factory()
	t := getElemType(obj)
	n := getShortName(t)
	all := append([]string{n}, names...)

	return &Entry{
		Names:        all,
		Type:         t,
		InternalType: internalType,
		Factory:      factory,
	}
}
