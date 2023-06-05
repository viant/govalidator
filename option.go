package govalidator

import "github.com/viant/structology"

type (
	Options struct {
		UseMarker bool
		marker    *structology.Marker
		Shallow   bool
		Path      *Path
	}
	Option func(c *Options)
)

//WithSetMarker creates with marker option
func WithSetMarker() Option {
	return func(c *Options) {
		c.UseMarker = true
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
