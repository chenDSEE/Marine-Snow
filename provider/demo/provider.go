package demo

import (
	"MarineSnow/framework"
)

const Key = "counter:provider"

func (csp *CounterServiceProvider) NewServiceProvider(container framework.ServiceContainer) framework.NewProviderFunc {
	return newCounterServiceProvider
}

func (csp *CounterServiceProvider) Init(container framework.ServiceContainer) error {
	return nil // do nothing
}

func (csp *CounterServiceProvider) IsDefer() bool {
	return false
}

func (csp *CounterServiceProvider) DefaultParams(container framework.ServiceContainer) []interface{} {
	return []interface{}{Key} // key as default
}

func (csp *CounterServiceProvider) Key() string {
	return Key
}
