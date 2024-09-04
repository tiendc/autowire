package autowire

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValueProvider(t *testing.T) {
	t.Run("Success with primitive type", func(t *testing.T) {
		p := newValueProvider(123, reflect.ValueOf(123))
		assert.Equal(t, 1, len(p.TargetTypes()))
		assert.Equal(t, reflect.TypeFor[int](), p.TargetTypes()[0])
		assert.Equal(t, 0, len(p.DependentTypes()))
	})

	t.Run("Success with pointer type", func(t *testing.T) {
		v := "abc"
		p := newValueProvider(&v, reflect.ValueOf(&v))
		assert.Equal(t, 1, len(p.TargetTypes()))
		assert.Equal(t, reflect.TypeFor[*string](), p.TargetTypes()[0])
		assert.Equal(t, 0, len(p.DependentTypes()))
	})
}
