package govalidator

type (
	Options struct {
		PresenceProvider *PresenceProvider
		Shallow          bool
		Path             *Path
	}
	Option func(c *Options)
)

//WithPresence creates with presence option
func WithPresence() Option {
	return func(c *Options) {
		c.PresenceProvider = &PresenceProvider{}
	}
}

//WithPath creates with path option
func WithPath(path *Path) Option {
	return func(c *Options) {
		c.Path = path
	}
}

//WithShallow creates with shallow option
func WithShallow(flag bool) Option {
	return func(c *Options) {
		c.Shallow = flag
	}
}

//newOptions creates an options
func newOptions() *Options {
	return &Options{}
}
