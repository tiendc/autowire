package autowire

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainerResolve_Failure(t *testing.T) {
	t.Run("Provider not found (user type)", func(t *testing.T) {
		c, err := NewContainer([]any{NewSrv1_OK_With_Need_Srv2_Srv3, NewSrv2_OK})
		assert.Nil(t, err)
		_, err = Resolve[Service1](c)
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Contains(t, err.Error(),
			"ErrNotFound: provider not found for type 'autowire.Service3'")
	})

	t.Run("Provider not found (primitive type)", func(t *testing.T) {
		c, err := NewContainer([]any{NewSrv1_OK_With_Need_Srv2_Srv3_Struct1, NewSrv2_OK, NewSrv3_OK})
		assert.Nil(t, err)
		_, err = Resolve[Service1](c)
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Contains(t, err.Error(),
			"ErrNotFound: provider not found for type '*autowire.Struct1_OK'")
	})

	t.Run("Provider not found for target type", func(t *testing.T) {
		c, err := NewContainer([]any{NewSrv1_OK, NewSrv2_OK})
		assert.Nil(t, err)
		_, err = Resolve[Service3](c)
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Contains(t, err.Error(),
			"ErrNotFound: provider not found for type 'autowire.Service3'")
	})

	t.Run("Circular dependency", func(t *testing.T) {
		// S1 -> S2 -> S4 -> S1
		c, err := NewContainer([]any{NewSrv1_OK_With_Need_Srv2_Srv3, NewSrv2_OK_With_Need_Srv4_Srv5,
			NewSrv5_OK, NewSrv3_OK, NewSrv4_OK_With_Need_Srv1})
		assert.Nil(t, err)
		_, err = Resolve[Service1](c)
		assert.ErrorIs(t, err, ErrCircularDependency)
		assert.Contains(t, err.Error(),
			"ErrCircularDependency: circular dependency detected at type 'autowire.Service1'")
	})

	t.Run("Circular dependency (self dependent)", func(t *testing.T) {
		// S1 -> S1
		c, err := NewContainer([]any{NewSrv1_Fail_Need_Srv1})
		assert.Nil(t, err)
		_, err = Resolve[Service1](c)
		assert.ErrorIs(t, err, ErrCircularDependency)
		assert.Contains(t, err.Error(),
			"ErrCircularDependency: circular dependency detected at type 'autowire.Service1'")
	})

	t.Run("Requires context.Context, but not provide", func(t *testing.T) {
		c, err := NewContainer([]any{NewSrv1_OK_With_Need_Ctx})
		assert.Nil(t, err)
		_, err = Resolve[Service1](c)
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Contains(t, err.Error(),
			"ErrNotFound: provider not found for type 'context.Context'")
	})
}

func TestContainerResolve_Success(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		c, err := NewContainer([]any{NewSrv1_OK_With_Need_Srv2_Srv3_IntSlice, NewSrv2_OK_With_Need_Srv4_Srv5,
			NewSrv3_OK, NewSrv4_OK, NewSrv5_OK, &struct1_OK, &struct5_OK})
		assert.Nil(t, err)
		dg, err := c.Resolve(reflect.TypeFor[Service1]())
		assert.Nil(t, err)
		assert.Equal(t, reflect.TypeFor[Service1](), dg.TargetType)
		assert.Equal(t, 3, len(dg.Dependencies))

		dg1 := dg.Dependencies[0]
		assert.Equal(t, reflect.TypeFor[Service2](), dg1.TargetType)
		assert.Equal(t, 2, len(dg1.Dependencies))
		dg11 := dg1.Dependencies[0]
		assert.Equal(t, reflect.TypeFor[Service4](), dg11.TargetType)
		assert.Equal(t, 0, len(dg11.Dependencies))
		dg12 := dg1.Dependencies[1]
		assert.Equal(t, reflect.TypeFor[Service5](), dg12.TargetType)
		assert.Equal(t, 0, len(dg12.Dependencies))

		dg2 := dg.Dependencies[1]
		assert.Equal(t, reflect.TypeFor[Service3](), dg2.TargetType)
		assert.Equal(t, 0, len(dg2.Dependencies))

		dg3 := dg.Dependencies[2]
		assert.Equal(t, reflect.TypeFor[[]int](), dg3.TargetType)
		assert.Equal(t, 0, len(dg3.Dependencies))
	})
}
