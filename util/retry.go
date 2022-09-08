package util

import (
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

// 重試 無限
func (r *Retry) RunLoop(f func() (error, interface{})) (error, interface{}) {
	var err error
	var res interface{}

	// 先跑第一次，若成功則直接回傳
	err, res = f()
	if err == nil {
		return nil, res
	}

	tick := time.NewTicker(r.interval * time.Millisecond)
	defer tick.Stop()

	for range tick.C {
		err, res := f()

		if err == nil {
			return nil, res
		}
	}

	return err, res
}

// 重試 限制次數
func (r *Retry) RunCount(f func() (error, interface{})) (error, interface{}) {
	var err error
	var res interface{}

	// 先跑第一次，若成功則直接回傳
	err, res = f()
	if err == nil {
		return nil, res
	}

	tick := time.NewTicker(r.interval * time.Millisecond)
	defer tick.Stop()

	maxCount := r.count

	for range tick.C {
		err, res := f()

		if err == nil {
			return nil, res
		}

		maxCount--
		if maxCount <= 0 {
			return err, nil
		}
	}

	return err, res
}
