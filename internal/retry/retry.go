package retry

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	// default params for retry
	defaultMaxRetry  = 5
	defaultBaseDelay = 1 * time.Second
	defaultMaxDelay  = 30 * time.Second
)

type RetryConfig struct {
	MaxRetry  uint
	BaseDelay time.Duration
	MaxDelay  time.Duration
	UsedFrom  string
}

func Retry(ctx context.Context, cfg RetryConfig, action func() error) error {
	validate(&cfg)
	var err error

	for attempt := range cfg.MaxRetry {
		err = action()
		if err == nil {
			return nil
		}

		backoff := cfg.BaseDelay * time.Duration(math.Pow(2, float64(attempt)))
		if backoff > cfg.MaxDelay {
			backoff = cfg.MaxDelay
		}

		jitter := time.Duration(rand.Int63n(int64(backoff / 2)))

		fmt.Printf(
			"%s attempt %d failed, waiting %v + jitter time %v (max delay %v)...\n",
			cfg.UsedFrom, attempt+1, backoff, jitter, cfg.MaxDelay,
		)
		time.Sleep(backoff + jitter)
	}

	return err
}

func validate(cfg *RetryConfig) {
	if cfg.MaxRetry == 0 {
		cfg.MaxRetry = defaultMaxRetry
	}
	if cfg.MaxDelay == 0 {
		cfg.MaxDelay = defaultMaxDelay
	}
	if cfg.BaseDelay == 0 {
		cfg.BaseDelay = defaultBaseDelay
	}
}
