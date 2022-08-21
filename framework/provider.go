package framework

type NewProviderFunc func(...interface{}) (interface{}, error)

// ServiceProvider interface specifies the behavior that three-part library must have
// to register into MarineSnow framework and created by MarineSnow framework.
// All method will be called by ServiceContainer
type ServiceProvider interface {
	// called by ServiceContainer to register ServiceProvider's NewProviderFunc
	NewServiceProvider(container ServiceContainer) NewProviderFunc

	// to setup ServiceProvider instance if need
	Init(container ServiceContainer) error

	// ture, initialize ServiceProvider intance when register
	// false, initialize ServiceProvider instance when user call ServiceContainer.Make()
	IsDefer() bool

	// set default params pass to NewProviderFunc
	DefaultParams(container ServiceContainer) []interface{}

	// get ServiceProvider name
	Key() string
}
