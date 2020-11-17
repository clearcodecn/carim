package consul

import (
	"context"
	"time"
)

var (
	DefaultOption = Options{
		Address: "http://127.0.0.1:8500",
		Timeout: 30 * time.Second,
	}
)

type Option func(*Options)

type Options struct {
	Address string // consul address
	Timeout time.Duration
	TTL     time.Duration
}

func Address(addr string) Option {
	return func(options *Options) {
		options.Address = addr
	}
}

func TTL(d time.Duration) Option {
	return func(options *Options) {
		options.TTL = d
	}
}

var (
	optionKey struct{}
)

func RegistryOption(opts ...Option) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, optionKey, opts)
	return ctx
}

func OptionFromContext(ctx context.Context) []Option {
	opts, _ := ctx.Value(optionKey).([]Option)
	return opts
}
