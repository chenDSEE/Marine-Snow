package demo

import (
	"MarineSnow/framework"
	"fmt"
)

const Key = "counter:provider"

func (csp *CounterServiceProvider) NewServiceProvider(container framework.ServiceContainer) framework.NewProviderFunc {
	return newCounterServiceProvider
}

func (csp *CounterServiceProvider) Init(container framework.ServiceContainer, params ...interface{}) error {
	// just log when CounterServiceProvider creating
	fmt.Println("create CounterServiceProvider with params:", params)
	return nil
}

func (csp *CounterServiceProvider) IsDefer() bool {
	//return false
	return true
}

func (csp *CounterServiceProvider) DefaultParams(container framework.ServiceContainer) []interface{} {
	return []interface{}{Key} // key as default
}

func (csp *CounterServiceProvider) Key() string {
	return Key
}
