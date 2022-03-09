package types

// TypeFunc func stuff
type TypeFunc func() interface{}

type Entry struct {
	Names   []string
	Factory TypeFunc
}

func NewEntryFromObject(obj interface{}, previousNames ...string) *Entry {
	t := GetElemType(obj)
	name := GetShortName(t)
	factory := TypeFactory(t)
	return NewEntry(name, factory, previousNames...)
}

func NewEntry(name string, factory TypeFunc, previousNames ...string) *Entry {
	all := append([]string{name}, previousNames...)
	return &Entry{
		Names:   all,
		Factory: factory,
	}
}
