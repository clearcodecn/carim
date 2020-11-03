package network

import "time"

type Identify interface {
	ID() string

	Token() string

	Expire() time.Time
}
