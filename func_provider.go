package autowire

import (
	"fmt"
	"reflect"
)

// funcProvider can take a function and execute it to create the target object
type funcProvider struct {
	baseProvider
}

// TargetTypes implementation of Provider interface. Typically, this returns
// a slice of one item which is the first return type of the function.
func (p *funcProvider) TargetTypes() []reflect.Type {
	return []reflect.Type{p.sourceVal.Type().Out(0)}
}

// DependentTypes implementation of Provider interface.
// This returns a slice of all types of the input arguments of the function.
func (p *funcProvider) DependentTypes() []reflect.Type {
	typ := p.sourceVal.Type()
	numIn := typ.NumIn()
	ret := make([]reflect.Type, 0, numIn)
	for i := 0; i < numIn; i++ {
		ret = append(ret, typ.In(i))
	}
	return ret
}

// parse parses and validates the source function
func (p *funcProvider) parse() error {
	typ := p.sourceVal.Type()
	if typ.Kind() != reflect.Func {
		return fmt.Errorf("%w: function required, got '%v'", ErrProviderInvalid, typ)
	}
	if typ.IsVariadic() {
		return fmt.Errorf("%w: variadic function is not allowed, error at '%v'", ErrProviderInvalid, typ)
	}

	// Validate function output
	numOut := typ.NumOut()
	if numOut == 0 || numOut > 2 {
		return fmt.Errorf(
			"%w: function must return either a <value> or a pair of (<value>, <error>), error at '%v'",
			ErrProviderInvalid, typ)
	}
	if numOut == 2 && !typ.Out(1).Implements(typeError) {
		return fmt.Errorf("%w: function must return second value of type error, got '%v', error at '%v'",
			ErrProviderInvalid, typ.Out(1), typ)
	}

	// Validate function input
	numIn := typ.NumIn()
	inArgTypes := make(map[reflect.Type]struct{}, numIn)
	for i := 0; i < numIn; i++ {
		inArg := typ.In(i)
		if _, exist := inArgTypes[inArg]; exist {
			return fmt.Errorf("%w: duplicated function argument type '%v', error at '%v'",
				ErrProviderInvalid, inArg, typ)
		}
		inArgTypes[inArg] = struct{}{}
	}

	return nil
}

// Build executes the source function and returns result.
// In case there are dependencies, this function will execute corresponding providers to
// collect all required objects to feed the current function.
func (p *funcProvider) Build(ctx *Context, targetType reflect.Type) (reflect.Value, error) {
	if ctx.sharedMode {
		if value, exist := ctx.objectMap[targetType]; exist {
			return value, nil
		}
	}

	if _, exist := ctx.resolvingTypes[targetType]; exist {
		return reflect.Value{}, fmt.Errorf("%w: circular dependency detected at type '%v'",
			ErrCircularDependency, targetType)
	}
	ctx.resolvingTypes[targetType] = struct{}{}
	defer func() {
		delete(ctx.resolvingTypes, targetType)
	}()

	var inArgs []reflect.Value
	for _, dependentType := range p.DependentTypes() {
		argProv, err := ctx.providerSet.GetFor(dependentType)
		if err != nil {
			return reflect.Value{}, err
		}
		argVal, err := argProv.Build(ctx, dependentType)
		if err != nil {
			return reflect.Value{}, err
		}
		inArgs = append(inArgs, argVal)
	}

	result := p.sourceVal.Call(inArgs)
	var err error
	if len(result) == 2 && result[1].IsValid() {
		iface := result[1].Interface()
		if iface != nil {
			err, _ = iface.(error)
		}
	}

	if err == nil && ctx.sharedMode {
		ctx.objectMap[targetType] = result[0]
	}
	return result[0], err
}

// newFuncProvider create a function provider
func newFuncProvider(provSrc any, provVal reflect.Value) (*funcProvider, error) {
	provider := &funcProvider{
		baseProvider: baseProvider{source: provSrc, sourceVal: provVal},
	}
	if err := provider.parse(); err != nil {
		return nil, err
	}
	return provider, nil
}
