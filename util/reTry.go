package util

import (
	"context"
	"time"

	"github.com/Chu16537/gomodules/env"
)

// 重試
func Retry(c context.Context, fn func() error) {
	tick := time.NewTicker(env.Env.RetryTime * time.Millisecond)
	defer tick.Stop()

	select {
	case <-c.Done():
		return
	case <-tick.C:
		err := fn
		if err == nil {
			return
		}
	}

}
