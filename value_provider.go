package autowire

import (
	"reflect"
)

// valueProvider can take a value of a certain type and return it when called
type valueProvider struct {
	baseProvider
	targetType reflect.Type
}

// TargetTypes returns a slice of target types
func (p *valueProvider) TargetTypes() []reflect.Type {
	return []reflect.Type{p.targetType}
}

// DependentTypes returns a slice of dependent types
func (p *valueProvider) DependentTypes() []reflect.Type {
	return []reflect.Type{}
}

// Build returns the value hold by this provider
func (p *valueProvider) Build(ctx *Context, targetType reflect.Type) (reflect.Value, error) {
	return p.sourceVal, nil
}

// newValueProvider creates a value provider
func newValueProvider[T any](provSrc T, provVal reflect.Value) *valueProvider {
	return &valueProvider{
		baseProvider: baseProvider{source: provSrc, sourceVal: provVal},
		targetType:   typeFor[T](),
	}
}
