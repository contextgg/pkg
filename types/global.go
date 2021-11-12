package types

var g = NewRegistry()

func Add(options ...TypeOption) (*Entry, error) {
	return g.Add(options...)
}

func GetByName(name string) (*Entry, bool) {
	return g.GetByName(name)
}
