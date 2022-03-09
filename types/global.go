package types

var g = NewRegistry()

func Add(e *Entry) error {
	return g.Add(e)
}

func GetByName(name string) (*Entry, bool) {
	return g.GetByName(name)
}
