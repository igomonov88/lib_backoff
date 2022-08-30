package backoff

import "time"

type Config struct {
	Min     time.Duration `validate:"nonzero"`
	Max     time.Duration `validate:"nonzero"`
	Factor  float64       `validate:"nonzero"`
	Retries int           `validate:"nonzero"`
}
