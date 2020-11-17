package server

import (
	"crypto/tls"
	"github.com/clearcodecn/carim/registry"
	"github.com/google/uuid"
)

type Option func(*Options)

/**
type Options struct {
	Codecs       map[string]codec.NewCodec
	Broker       broker.Broker
	Registry     registry.Registry
	Tracer       trace.Tracer
	Auth         auth.Auth
	Transport    transport.Transport
	Metadata     map[string]string
	Name         string
	Address      string
	Advertise    string
	Id           string
	Version      string
	HdlrWrappers []HandlerWrapper
	SubWrappers  []SubscriberWrapper

	// RegisterCheck runs a check function before registering the service
	RegisterCheck func(context.Context) error
	// The register expiry time
	RegisterTTL time.Duration
	// The interval on which to register
	RegisterInterval time.Duration

	// The router for requests
	Router Router

	// TLSConfig specifies tls.Config for secure serving
	TLSConfig *tls.Config

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}
*/

type Options struct {
	Id             string
	Version        string
	Address        string
	TLSConfig      *tls.Config
	RegistryOption registry.Option

	Transport Transport
}

func Id(id string) Option {
	return func(options *Options) {
		options.Id = id
	}
}

func Version(version string) Option {
	return func(options *Options) {
		options.Version = version
	}
}

func Address(address string) Option {
	return func(options *Options) {
		options.Address = address
	}
}

func WithTLS(config *tls.Config) Option {
	return func(options *Options) {
		options.TLSConfig = config
	}
}

var (
	DefaultOption Options
)

func init() {
	DefaultOption = Options{
		Id:             uuid.New().String(),
		Address:        ":9000",
		TLSConfig:      nil,
		RegistryOption: registry.DefaultOption,
	}
}
