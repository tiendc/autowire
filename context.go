package autowire

import (
	"reflect"
)

// Context a context object used in each building/resolving object
type Context struct {
	sharedMode bool

	providerSet ProviderSet
	objectMap   map[reflect.Type]reflect.Value

	resolvingTypes map[reflect.Type]struct{}
}

// ContextOption configuration setter for a context
type ContextOption func(*Context)

// NonSharedMode sets `shared mode` to `false` for the current context
func NonSharedMode() ContextOption {
	return func(ctx *Context) {
		ctx.sharedMode = false
	}
}

// ProviderOverwrite overwrites a value for the current context
func ProviderOverwrite[T any](val T) ContextOption {
	return func(ctx *Context) {
		ctx.providerSet.Overwrite(newValueProvider(val, reflect.ValueOf(val)))
	}
}
