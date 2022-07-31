package es

func dataOptions(options []DataOption) DataOptions {
	// set defaults.
	o := DataOptions{}
	// apply options.
	for _, opt := range options {
		opt(&o)
	}
	return o
}

type DataOption func(d *DataOptions)

type DataOptions struct {
	RecreateTables bool
	TruncateTables bool
	HasEvents      bool
	HasSnapshots   bool
	ExtraModels    []interface{}
}

func InitializeSnapshots() DataOption {
	return func(o *DataOptions) {
		o.HasSnapshots = true
	}
}
func InitializeEvents() DataOption {
	return func(o *DataOptions) {
		o.HasEvents = true
	}
}
func InitializeEntities(entities ...Entity) DataOption {
	return func(o *DataOptions) {
		for _, ent := range entities {
			o.ExtraModels = append(o.ExtraModels, ent)
		}
	}
}
func WithTruncate() DataOption {
	return func(o *DataOptions) {
		o.TruncateTables = true
	}
}
func WithRecreate() DataOption {
	return func(o *DataOptions) {
		o.RecreateTables = true
	}
}
