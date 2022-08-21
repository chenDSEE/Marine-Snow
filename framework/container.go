package framework

import (
	"errors"
	"sync"
)

// TODO: support multi-ServiceProvider for same type
// ServiceContainer interface specifies the behavior that MarineSnow framework MUST have.
// The framework user can use below method to extract any provider in handler function.
type ServiceContainer interface {
	// Register() will replace the old provider with same ServiceProvider.Key()
	Register(provider ServiceProvider) error

	IsRegistered(key string) bool // TODO: key should be a special type, like: framework.ServiceProviderKey(type string)

	// get a ServiceProvider instance from container
	// instance maybe reuse
	Make(key string) (interface{}, error)

	// get a ServiceProvider instance from container, painc when error
	// instance maybe reuse
	MustMake(key string) interface{}

	// make new ServiceProvider instance with params, never store the instance into ServiceContainer
	MakeNew(key string, params ...interface{}) (interface{}, error)
}

type Container struct {
	sync.RWMutex

	providers map[string]ServiceProvider
	instances map[string]interface{}
}

func NewContainer() ServiceContainer {
	return &Container{
		providers: map[string]ServiceProvider{},
		instances: map[string]interface{}{},
	}
}

var _ ServiceContainer = &Container{}

func (container *Container) Register(provider ServiceProvider) error {
	container.Lock()
	defer container.Unlock()

	/* record provider */
	container.providers[provider.Key()] = provider

	/* create and init provider instance if need */
	if provider.IsDefer() == false {
		if err := provider.Init(container, provider.DefaultParams(container)...); err != nil {
			return err
		}

		newFunc := provider.NewServiceProvider(container)
		instance, err := newFunc(provider.DefaultParams(container))
		if err != nil {
			return err
		}

		container.instances[provider.Key()] = instance
	}

	return nil
}

func (container *Container) IsRegistered(key string) bool {
	container.RLock()
	defer container.RUnlock()

	_, ok := container.providers[key]
	return ok
}

func (container *Container) Make(key string) (interface{}, error) {
	return container.make(key, nil, false)
}

func (container *Container) MustMake(key string) interface{} {
	instance, err := container.make(key, nil, false)
	if err != nil {
		panic("fail to make instance for " + key + " ServiceProvider, err: " + err.Error())
	}

	return instance
}

// TODO: MakeNew() instance can be reused will be better
func (container *Container) MakeNew(key string, params ...interface{}) (interface{}, error) {
	return container.make(key, params, true)
}

// TODO: lock-free to make will be better, but provider can't update when server is running
func (container *Container) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	/* if provider existed */
	container.RLock()
	provider, ok := container.providers[key]
	if !ok {
		container.RUnlock()
		return nil, errors.New("SerivceProvider with " + key + " not register into container")
	}

	if forceNew {
		container.RUnlock()
		return container.newInstance(provider, params)
	}

	/* return existed provider instance */
	if instance, existed := container.instances[key]; existed {
		container.RUnlock()
		return instance, nil
	}

	container.RUnlock()

	/* create new provider instance and store it */
	container.Lock()
	defer container.Unlock()

	if instance, existed := container.instances[key]; existed {
		// NOTE: should double check after write lock, avoid concurrency problem
		return instance, nil
	}

	instance, err := container.newInstance(provider, params)
	if err != nil {
		return nil, err
	}

	container.instances[key] = instance
	return instance, nil
}

// create new instance, but not store into container
func (container *Container) newInstance(provider ServiceProvider, params []interface{}) (interface{}, error) {
	if params == nil {
		params = provider.DefaultParams(container)
	}

	if err := provider.Init(container, params...); err != nil {
		return nil, err
	}

	newFunc := provider.NewServiceProvider(container)
	instance, err := newFunc(params...)
	if err != nil {
		return nil, err
	}

	return instance, nil
}
