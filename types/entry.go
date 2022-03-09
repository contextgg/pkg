package types

// TypeFunc func stuff
type TypeFunc func() interface{}

type Entry struct {
	Name          string
	PreviousNames []string
	Factory       TypeFunc
}
