package types

var g = NewRegistry()

func Add(entries ...*Entry) error {
	return g.Add(entries...)
}

func GetByName(name string) (*Entry, bool) {
	return g.GetByName(name)
}
