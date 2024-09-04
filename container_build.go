package autowire

import (
	"context"
	"reflect"
)

// Build implementation of Container interface
func (c *container) Build(targetType reflect.Type, opts ...ContextOption) (value reflect.Value, err error) {
	provider, err := c.providerSet.GetFor(targetType)
	if err != nil {
		return value, err
	}

	ctx := &Context{
		sharedMode:     c.sharedMode,
		providerSet:    c.providerSet.shallowClone(),
		objectMap:      c.objectMap,
		resolvingTypes: make(map[reflect.Type]struct{}, 10), //nolint:gomnd
	}
	for _, opt := range opts {
		opt(ctx)
	}

	value, err = provider.Build(ctx, targetType)
	if err != nil {
		return value, err
	}

	return value, nil
}

// BuildWithCtx implementation of Container interface
func (c *container) BuildWithCtx(ctx context.Context, targetType reflect.Type, opts ...ContextOption) (
	value reflect.Value, err error,
) {
	return c.Build(targetType, append(opts, ProviderOverwrite(ctx))...)
}
