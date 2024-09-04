package autowire

import (
	"fmt"
	"reflect"
)

// DependencyGraph dependency graph info of a type.
type DependencyGraph struct {
	TargetType   reflect.Type
	Dependencies []DependencyGraph
}

// Resolve implementation of Container interface
func (c *container) Resolve(targetType reflect.Type) (value DependencyGraph, err error) {
	ctx := &Context{
		sharedMode:     c.sharedMode,
		providerSet:    c.providerSet,
		objectMap:      c.objectMap,
		resolvingTypes: make(map[reflect.Type]struct{}, 10), //nolint:gomnd
	}
	return c.resolve(ctx, targetType)
}

func (c *container) resolve(ctx *Context, targetType reflect.Type) (DependencyGraph, error) {
	if _, exist := ctx.resolvingTypes[targetType]; exist {
		return DependencyGraph{}, fmt.Errorf("%w: circular dependency detected at type '%v'",
			ErrCircularDependency, targetType)
	}
	ctx.resolvingTypes[targetType] = struct{}{}

	provider, err := ctx.providerSet.GetFor(targetType)
	if err != nil {
		return DependencyGraph{}, err
	}

	depGraph := DependencyGraph{
		TargetType: targetType,
	}
	for _, dType := range provider.DependentTypes() {
		dGraph, err := c.resolve(ctx, dType)
		if err != nil {
			return DependencyGraph{}, err
		}
		depGraph.Dependencies = append(depGraph.Dependencies, dGraph)
	}

	delete(ctx.resolvingTypes, targetType)
	return depGraph, nil
}
