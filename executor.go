package backoff

import (
	"context"
	"fmt"
	"time"

	"github.com/jpillora/backoff"
)

// Execute run `fn` function in backoff flow
func Execute(fn func() error, cfg Config) (err error) {
	b := backoff.Backoff{
		Min:    cfg.Min,
		Max:    cfg.Max,
		Factor: cfg.Factor,
	}

	for b.Attempt() < float64(cfg.Retries) {
		err = fn()

		switch err.(type) {
		case ErrCancel:
			return
		}

		switch err {
		case nil:
		case context.Canceled:
		case context.DeadlineExceeded:
		default:
			time.Sleep(b.Duration())
			continue
		}

		return
	}

	return fmt.Errorf("backoff failed after %v attempts: %v", cfg.Retries, err)
}

type ErrCancel struct {
	error
}

// NewErrCancel creates an error that stops backoff execution
func NewErrCancel(e error) ErrCancel {
	return ErrCancel{fmt.Errorf("backoff was stopped manually, err: %s", e.Error())}
}
