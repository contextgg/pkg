package types

func GetFirstByNames(r Registry, names []string) (*Entry, bool) {
	for _, name := range names {
		if entry, ok := r.GetByName(name); ok {
			return entry, true
		}
	}
	return nil, false
}
