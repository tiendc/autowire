package autowire

import (
	"fmt"
	"reflect"
)

// ProviderSet is a set of unique providers for certain unique types
type ProviderSet interface {
	// GetFor returns a provider for the specified type, or ErrNotFound.
	GetFor(reflect.Type) (Provider, error)
	// GetAll returns all providers contained within the set
	GetAll() []Provider
	// Overwrite replaces the existing provider by its target type with the specified one
	Overwrite(Provider)
	// shallowClone clones the set (shallow clone only)
	shallowClone() ProviderSet
}

// providerSet default implementation of ProviderSet interface
type providerSet struct {
	providerMap            map[reflect.Type]Provider
	overwrittenProviderMap map[reflect.Type]Provider
}

// GetFor implementation of ProviderSet interface
func (ps *providerSet) GetFor(targetType reflect.Type) (Provider, error) {
	if len(ps.overwrittenProviderMap) > 0 {
		if prov, exist := ps.overwrittenProviderMap[targetType]; exist {
			return prov, nil
		}
	}
	if prov, exist := ps.providerMap[targetType]; exist {
		return prov, nil
	}
	return nil, fmt.Errorf("%w: provider not found for type '%v'", ErrNotFound, targetType)
}

// GetAll implementation of ProviderSet interface
func (ps *providerSet) GetAll() []Provider {
	numOverwritten := len(ps.overwrittenProviderMap)
	ret := make([]Provider, 0, len(ps.providerMap)+numOverwritten)
	for typ, v := range ps.providerMap {
		if numOverwritten > 0 {
			if _, exist := ps.overwrittenProviderMap[typ]; exist {
				continue
			}
		}
		ret = append(ret, v)
	}
	for _, v := range ps.overwrittenProviderMap {
		ret = append(ret, v)
	}
	return ret
}

// Overwrite implementation of ProviderSet interface
func (ps *providerSet) Overwrite(provider Provider) {
	if ps.overwrittenProviderMap == nil {
		ps.overwrittenProviderMap = map[reflect.Type]Provider{}
	}
	ps.overwrittenProviderMap[provider.TargetTypes()[0]] = provider
}

// shallowClone implementation of ProviderSet interface
func (ps *providerSet) shallowClone() ProviderSet {
	return &providerSet{
		providerMap:            ps.providerMap,
		overwrittenProviderMap: ps.overwrittenProviderMap,
	}
}

// NewProviderSet creates a new provider set from individual providers.
// A provider object can be:
//   - function in the below forms:
//     func(<arg1>, ..., <argN>) <ServiceType>
//     func(<arg1>, ..., <argN>) (<ServiceType>, error)
//     func(context.Context, <arg1>, ..., <argN>) (<ServiceType>, error)
//   - a struct pointer
//   - an object of type `ProviderSet`
//   - an object of type `Provider`
func NewProviderSet(args ...any) (ProviderSet, error) {
	return parseProviders(args...)
}

// MustNewProviderSet creates a ProviderSet with panicking on error
func MustNewProviderSet(args ...any) ProviderSet {
	ps, err := NewProviderSet(args...)
	if err != nil {
		panic(err)
	}
	return ps
}

//nolint:gocognit
func parseProviders(args ...any) (ProviderSet, error) {
	providerMap := make(map[reflect.Type]Provider, len(args))
	var err error

	for _, provSrc := range args {
		if provSrc == nil {
			return nil, fmt.Errorf("%w: provider must not be nil", ErrProviderInvalid)
		}

		if srcIsProvSet, ok := provSrc.(ProviderSet); ok {
			for _, prov := range srcIsProvSet.GetAll() {
				if err = addProviderToMap(prov, providerMap); err != nil {
					return nil, err
				}
			}
			continue
		}
		if srcIsProv, ok := provSrc.(Provider); ok {
			if err = addProviderToMap(srcIsProv, providerMap); err != nil {
				return nil, err
			}
			continue
		}

		provVal := reflect.ValueOf(provSrc)
		kind := provVal.Kind()
		if kind == reflect.Interface || kind == reflect.Pointer {
			provVal = provVal.Elem()
			if !provVal.IsValid() {
				return nil, fmt.Errorf("%w: provider must not be nil", ErrProviderInvalid)
			}
		}

		var provider Provider
		switch provVal.Kind() { //nolint:exhaustive
		case reflect.Func:
			provider, err = newFuncProvider(provSrc, provVal)
		case reflect.Struct:
			if kind != reflect.Pointer {
				return nil, fmt.Errorf("%w: struct pointer required, got '%v'",
					ErrProviderInvalid, reflect.TypeOf(provSrc))
			}
			provider, err = newStructProvider(provSrc, provVal)
		default:
			return nil, fmt.Errorf("%w: provider type unsupported, got '%v'",
				ErrProviderInvalid, reflect.TypeOf(provSrc))
		}
		if err != nil {
			return nil, err
		}
		if err = addProviderToMap(provider, providerMap); err != nil {
			return nil, err
		}
	}

	if len(providerMap) == 0 {
		return nil, fmt.Errorf("%w: no provider provided", ErrProviderInvalid)
	}
	return &providerSet{
		providerMap: providerMap,
	}, nil
}

func addProviderToMap(provider Provider, providerMap map[reflect.Type]Provider) error {
	for _, targetType := range provider.TargetTypes() {
		if _, exist := providerMap[targetType]; exist {
			return fmt.Errorf("%w: duplicated provider for type '%v'",
				ErrProviderDuplicated, targetType)
		}
		providerMap[targetType] = provider
	}
	return nil
}
