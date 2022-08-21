package demo

import (
	"errors"
	"sync/atomic"
)

type CounterServiceProvider struct {
	name string
	cnt  int64
}

func newCounterServiceProvider(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("error params numbers, it should be 1")
	}

	name, ok := params[0].(string)
	if !ok {
		return nil, errors.New("error params type, it should be string")
	}

	instance := &CounterServiceProvider{
		name: name,
	}

	return instance, nil
}

// CounterService interface specifies what the behavior CounterService should do
// all method should be concurrency safe
type CounterService interface {
	Name() string
	Cnt() int64

	// action
	Increase()
	Decrease()
	Reset()
}

func (csp *CounterServiceProvider) Cnt() int64 {
	return atomic.LoadInt64(&csp.cnt)
}

func (csp *CounterServiceProvider) Name() string {
	return csp.name
}

func (csp *CounterServiceProvider) Increase() {
	atomic.AddInt64(&csp.cnt, 1)
}

func (csp *CounterServiceProvider) Decrease() {
	atomic.AddInt64(&csp.cnt, -1)
}

func (csp *CounterServiceProvider) Reset() {
	atomic.StoreInt64(&csp.cnt, 0)
}
