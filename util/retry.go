package util

import (
	"context"
	"time"
)

type Retry struct {
	count    int           // 次數
	interval time.Duration // 間隔
}

// 取得
func NewRetry(count int, interval time.Duration) *Retry {
	return &Retry{
		count:    count,
		interval: interval,
	}

}

func (r *Retry) Run(c context.Context, fn func() (error, interface{})) (error, interface{}) {
	var err error
	var result interface{}

	tick := time.NewTicker(r.interval * time.Millisecond)
	defer tick.Stop()

	maxCount := r.count

	select {
	case <-c.Done():
		return nil, nil
	case <-tick.C:
		err, result = fn()
		if err == nil {
			return nil, result
		}

		maxCount--
		if maxCount <= 0 {
			return err, nil
		}
	}

	return err, result
}
