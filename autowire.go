package autowire

import (
	"context"
	"fmt"
	"reflect"
)

// Build builds object for the specified type within a container
func Build[T any](c Container, opts ...ContextOption) (value T, err error) {
	targetType := reflect.TypeFor[T]()

	val, err := c.Build(targetType, opts...)
	if err != nil {
		return value, err
	}

	value, ok := val.Interface().(T)
	if !ok { // this should never happen
		return value, fmt.Errorf("%w: unable to cast result as type '%v'", ErrTypeCast, targetType)
	}

	return value, nil
}

// BuildWithCtx builds object for the specified type within a container.
// This function will pass the specified context object to every provider that requires a context.
func BuildWithCtx[T any](ctx context.Context, c Container, opts ...ContextOption) (value T, err error) {
	targetType := reflect.TypeFor[T]()

	val, err := c.BuildWithCtx(ctx, targetType, opts...)
	if err != nil {
		return value, err
	}

	value, ok := val.Interface().(T)
	if !ok { // this should never happen
		return value, fmt.Errorf("%w: unable to cast result as type '%v'", ErrTypeCast, targetType)
	}

	return value, nil
}

// Get gets object of a type within a container.
// If no object is created for the type or `sharedMode` is `false`, ErrNotFound is returned.
func Get[T any](c Container) (value T, err error) {
	targetType := reflect.TypeFor[T]()

	val, err := c.Get(targetType)
	if err != nil {
		return value, err
	}

	value, ok := val.Interface().(T)
	if !ok { // this should never happen
		return value, fmt.Errorf("%w: unable to cast result as type '%v'", ErrTypeCast, targetType)
	}

	return value, nil
}

// Resolve builds dependency graph for the specified type within a container
func Resolve[T any](c Container) (DependencyGraph, error) {
	return c.Resolve(reflect.TypeFor[T]())
}
