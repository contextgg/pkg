package types

var g = NewRegistry()

func Upsert(obj interface{}, replaced ...string) *Entry {
	return g.Upsert(obj, replaced...)
}

func GetByName(name string) (*Entry, bool) {
	return g.GetByName(name)
}
