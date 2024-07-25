package govalidator

import "github.com/viant/structology"

type CanUseMarkerProvider func(v interface{}) bool

type (
	Options struct {
		UseMarker            bool
		PreservePointer      bool
		marker               *structology.Marker
		Shallow              bool
		Path                 *Path
		CanUseMarkerProvider CanUseMarkerProvider
	}

	Option func(c *Options)
)

// WithSetMarker creates with marker option
func WithSetMarker() Option {
	return func(c *Options) {
		c.UseMarker = true
	}
}

// WithPath creates with path option
func WithPath(path *Path) Option {
	return func(c *Options) {
		c.Path = path
	}
}

// WithShallow creates with shallow option
func WithShallow(flag bool) Option {
	return func(c *Options) {
		c.Shallow = flag
	}
}

// WithPreservePointer creates with preserve pointer option
func WithPreservePointer(flag bool) Option {
	return func(c *Options) {
		c.PreservePointer = flag
	}
}

// WithCanUseMarkerProvider creates with marker provider option
func WithCanUseMarkerProvider(provider CanUseMarkerProvider) Option {
	return func(c *Options) {
		c.UseMarker = true
		c.CanUseMarkerProvider = provider
	}
}

// newOptions creates an options
func newOptions() *Options {
	return &Options{}
}
