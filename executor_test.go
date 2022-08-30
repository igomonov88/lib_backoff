package backoff

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecuteWithError(t *testing.T) {
	var calls int32

	err := Execute(func() error {
		atomic.AddInt32(&calls, 1)
		return errors.New("something")
	}, Config{
		Min:     50 * time.Millisecond,
		Max:     200 * time.Millisecond,
		Factor:  1.5,
		Retries: 3,
	})

	assert.Equal(t, int32(3), calls, "it should be equal to Config.Retries")
	require.NotNil(t, err)
	t.Log(err.Error())
}

func TestExecuteCanceled(t *testing.T) {
	var calls int32

	ctx, cancel := context.WithCancel(context.Background())

	err := Execute(func() error {
		atomic.AddInt32(&calls, 1)
		cancel()
		return ctx.Err()
	}, Config{
		Min:     50 * time.Millisecond,
		Max:     200 * time.Millisecond,
		Factor:  1.5,
		Retries: 3,
	})

	assert.Equal(t, int32(1), calls, "it should be canceled")
	require.NotNil(t, err)
	t.Log(err.Error())
}

func TestExecuteDeadlineExceeded(t *testing.T) {
	var calls int32

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)

	err := Execute(func() error {
		atomic.AddInt32(&calls, 1)
		time.Sleep(100 * time.Millisecond)
		cancel()
		return ctx.Err()
	}, Config{
		Min:     50 * time.Millisecond,
		Max:     200 * time.Millisecond,
		Factor:  1.5,
		Retries: 3,
	})

	assert.Equal(t, int32(1), calls, "it should exceed the deadline")
	require.NotNil(t, err)
	t.Log(err.Error())
}

func TestExecuteErrCancel(t *testing.T) {
	var calls int32

	err := Execute(func() error {
		atomic.AddInt32(&calls, 1)
		return NewErrCancel(errors.New("something"))
	}, Config{
		Min:     50 * time.Millisecond,
		Max:     200 * time.Millisecond,
		Factor:  1.5,
		Retries: 3,
	})

	assert.Equal(t, int32(1), calls, "it should be equal to Config.Retries")
	require.NotNil(t, err)
	t.Log(err.Error())
}
