package es

// EntityOptions represents the configuration options
// for the entity.
type EntityOptions struct {
	Factory        EntityFunc
	Revision       string
	MinVersionDiff int
	Project        bool
}

// EntityOption applies an option to the provided configuration.
type EntityOption func(*EntityOptions)

func EntityRevision(revision string) EntityOption {
	return func(o *EntityOptions) {
		o.Revision = revision
	}
}
func EntityRevisionMin(minVersionDiff int) EntityOption {
	return func(o *EntityOptions) {
		o.MinVersionDiff = minVersionDiff
	}
}
func EntityDisableRevision() EntityOption {
	return func(o *EntityOptions) {
		o.MinVersionDiff = -1
	}
}
func EntityDisableProject() EntityOption {
	return func(o *EntityOptions) {
		o.Project = false
	}
}

// EntityFactory specifies the option to provide a factory for an entity.
func EntityFactory(factory EntityFunc) EntityOption {
	return func(o *EntityOptions) {
		o.Factory = factory
		o.Revision = "rev1"
		o.Project = true
	}
}
