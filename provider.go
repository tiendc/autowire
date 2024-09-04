package autowire

import (
	"reflect"
)

var (
	typeError = typeFor[error]()
)

// Provider a provider is an `object creator` and provide the object to a container.
// An object creation may require dependencies such as other objects or data.
// A `provider` should be able to resolve all dependencies and collect the required
// objects to create the expected object.
type Provider interface {
	// Source returns provider source which can be a function, a struct pointer, or a value
	Source() any
	// TargetTypes returns a list of target types that the provider can create objects of.
	// A function provider can provide only one target type, whereas a struct provider can
	// provide multiple target types through its fields.
	TargetTypes() []reflect.Type
	// DependentTypes returns a list of types that will be required when the provider creates
	// objects of the target types. Typically, a function provider can have no dependent types
	// or multiple ones, whereas a struct provider have no dependent types.
	DependentTypes() []reflect.Type
	// Build builds an object of the specified type.
	Build(*Context, reflect.Type) (reflect.Value, error)
}

// baseProvider base provider struct
type baseProvider struct {
	source    any
	sourceVal reflect.Value
}

// Source returns the provider source
func (p *baseProvider) Source() any {
	return p.source
}
