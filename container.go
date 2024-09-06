package autowire

import (
	"context"
	"fmt"
	"reflect"
)

// Container is a storage for storing every object created by the providers of the container.
// Each container has its own provider set and configuration settings.
type Container interface {
	// SharedMode gets shared mode in the container (default is `true`).
	//
	// When `shared mode` is `true`, every object created within the container will be cached and returned
	// on future uses. Therefore, there is only an object created for a specific type.
	// For example: ServiceA requires ServiceX and ServiceY. ServiceB requires ServiceX and ServiceZ.
	// In this mode, ServiceX will be created only one time, so ServiceA and ServiceB will share the
	// same ServiceX object.
	SharedMode() bool

	// setSharedMode sets shared mode
	setSharedMode(bool)

	// ProviderSet gets provider set within the container
	ProviderSet() ProviderSet

	// Get gets a value stored in the container for the specified type.
	// If not found, returns ErrNotFound.
	Get(targetType reflect.Type) (reflect.Value, error)

	// Build creates a value for the specified type and all other required values.
	Build(targetType reflect.Type, opts ...ContextOption) (reflect.Value, error)

	// BuildWithCtx creates a value for the specified type with passing a context.Context object.
	// The context object will be passed to every provider which requires a context.
	BuildWithCtx(ctx context.Context, targetType reflect.Type, opts ...ContextOption) (reflect.Value, error)

	// Resolve builds dependency graph for the specified type
	Resolve(targetType reflect.Type) (DependencyGraph, error)
}

// ContainerConfigOption config option setter used when create a container
type ContainerConfigOption func(Container)

// SetSharedMode config option for setting `sharedMode` for a container
func SetSharedMode(flag bool) ContainerConfigOption {
	return func(c Container) {
		c.setSharedMode(flag)
	}
}

// container an implementation of Container interface
type container struct {
	sharedMode  bool
	providerSet ProviderSet
	objectMap   map[reflect.Type]reflect.Value
}

// SharedMode implementation of Container interface
func (c *container) SharedMode() bool {
	return c.sharedMode
}

// setSharedMode implementation of Container interface
func (c *container) setSharedMode(flag bool) {
	c.sharedMode = flag
}

// ProviderSet implementation of Container interface
func (c *container) ProviderSet() ProviderSet {
	return c.providerSet
}

// Get implementation of Container interface
func (c *container) Get(targetType reflect.Type) (value reflect.Value, err error) {
	if c.objectMap != nil {
		if value, exist := c.objectMap[targetType]; exist {
			return value, nil
		}
	}
	return value, fmt.Errorf("%w: object not found for type '%v'", ErrNotFound, targetType)
}

// NewContainer creates a new container with providing providers and settings.
// Provider list can contain:
//   - functions in the below forms:
//     func(<arg1>, ..., <argN>) <ServiceType>
//     func(<arg1>, ..., <argN>) (<ServiceType>, error)
//     func(context.Context, <arg1>, ..., <argN>) (<ServiceType>, error)
//   - struct pointers
//   - objects of type `ProviderSet`
//   - objects of type `Provider`
func NewContainer(providers []any, opts ...ContainerConfigOption) (Container, error) {
	providerSet, err := parseProviders(providers...)
	if err != nil {
		return nil, err
	}

	c := &container{
		sharedMode:  true,
		providerSet: providerSet,
		objectMap:   map[reflect.Type]reflect.Value{},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// MustNewContainer creates a new container and panics on error.
func MustNewContainer(providers []any, opts ...ContainerConfigOption) Container {
	c, err := NewContainer(providers, opts...)
	if err != nil {
		panic(err)
	}
	return c
}
