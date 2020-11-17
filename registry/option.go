package registry

import "time"

type Option func(*Options)

type Options struct {
	TTL time.Duration


}
