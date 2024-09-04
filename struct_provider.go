package autowire

import (
	"fmt"
	"reflect"
	"strings"
)

// structProvider can take a struct pointer and return its fields based on specific type
type structProvider struct {
	baseProvider
	targetTypes map[reflect.Type]structFieldDetail
}

type structFieldDetail struct {
	Index []int
	Name  []string
}

// TargetTypes returns a slice of types in which the provider can return objects of
func (p *structProvider) TargetTypes() []reflect.Type {
	types := make([]reflect.Type, 0, len(p.targetTypes))
	for typ := range p.targetTypes {
		types = append(types, typ)
	}
	return types
}

// DependentTypes returns a slice of dependent types. This returns an empty slice as
// struct fields require no dependency.
func (p *structProvider) DependentTypes() []reflect.Type {
	return []reflect.Type{}
}

// parse parses the input struct to collect target types
func (p *structProvider) parse() error {
	typ := p.sourceVal.Type()
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("%w: struct required, got '%v'", ErrProviderInvalid, typ)
	}

	p.targetTypes = make(map[reflect.Type]structFieldDetail, 10) //nolint:gomnd
	if err := p.parseStructRecursively(typ, typ, []int{}, []string{}); err != nil {
		return err
	}
	return nil
}

func (p *structProvider) parseStructRecursively(rootType, typ reflect.Type, index []int, name []string) error {
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil
	}

	numField := typ.NumField()
	for i := 0; i < numField; i++ {
		sf := typ.Field(i)
		if !sf.IsExported() {
			continue
		}
		if detail, exist := p.targetTypes[sf.Type]; exist {
			path1 := strings.Join(detail.Name, ".")
			path2 := strings.Join(append(name, sf.Name), ".")
			return fmt.Errorf("%w: duplicated provider for type '%v', error at '%v[%s]' and '%v[%s]'",
				ErrProviderDuplicated, sf.Type, rootType, path1, rootType, path2)
		}
		detail := structFieldDetail{
			Index: append(index, i),
			Name:  append(name, sf.Name),
		}
		p.targetTypes[sf.Type] = detail
		if err := p.parseStructRecursively(rootType, sf.Type, detail.Index, detail.Name); err != nil {
			return err
		}
	}
	return nil
}

// Build returns related field value for the specified type
func (p *structProvider) Build(ctx *Context, targetType reflect.Type) (reflect.Value, error) {
	if fieldDetail, exist := p.targetTypes[targetType]; exist {
		val, err := p.sourceVal.FieldByIndexErr(fieldDetail.Index)
		if err != nil {
			return val, err
		}
		return val, nil
	}
	return reflect.Value{}, fmt.Errorf("%w: provider not found for '%v'", ErrNotFound, targetType)
}

// newStructProvider creates a struct provider
func newStructProvider(provSrc any, provVal reflect.Value) (*structProvider, error) {
	provider := &structProvider{
		baseProvider: baseProvider{source: provSrc, sourceVal: provVal},
	}
	if err := provider.parse(); err != nil {
		return nil, err
	}
	return provider, nil
}
