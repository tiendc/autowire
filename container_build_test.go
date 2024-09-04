package autowire

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainerBuild_Failure(t *testing.T) {
	t.Run("No provider provided", func(t *testing.T) {
		_, err := NewContainer(nil)
		assert.ErrorIs(t, err, ErrProviderInvalid)
		assert.Contains(t, err.Error(), "ErrProviderInvalid: no provider provided")
	})

	t.Run("Provider not found (user type)", func(t *testing.T) {
		c, err := NewContainer([]any{NewSrv1_OK_With_Need_Srv2_Srv3, NewSrv2_OK})
		assert.Nil(t, err)
		_, err = c.Build(reflect.TypeFor[Service1]())
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Contains(t, err.Error(),
			"ErrNotFound: provider not found for type 'autowire.Service3'")
	})

	t.Run("Provider not found (primitive type)", func(t *testing.T) {
		c, err := NewContainer([]any{NewSrv1_OK_With_Need_Srv2_Srv3_Struct1, NewSrv2_OK, NewSrv3_OK})
		assert.Nil(t, err)
		_, err = c.Build(reflect.TypeFor[Service1]())
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Contains(t, err.Error(),
			"ErrNotFound: provider not found for type '*autowire.Struct1_OK'")
	})

	t.Run("Provider not found for target type", func(t *testing.T) {
		c, err := NewContainer([]any{NewSrv1_OK, NewSrv2_OK})
		assert.Nil(t, err)
		_, err = c.Build(reflect.TypeFor[Service3]())
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Contains(t, err.Error(),
			"ErrNotFound: provider not found for type 'autowire.Service3'")
	})

	t.Run("Provider return error as 2nd value", func(t *testing.T) {
		c, err := NewContainer([]any{NewSrv1_Fail_With_Err, NewSrv2_OK, NewSrv3_OK})
		assert.Nil(t, err)
		_, err = c.Build(reflect.TypeFor[Service1]())
		assert.ErrorIs(t, err, errTest1)
		assert.Contains(t, err.Error(), "errTest1")
	})

	t.Run("Circular dependency", func(t *testing.T) {
		// S1 -> S2 -> S4 -> S1
		c, err := NewContainer([]any{NewSrv1_OK_With_Need_Srv2_Srv3, NewSrv2_OK_With_Need_Srv4_Srv5,
			NewSrv5_OK, NewSrv3_OK, NewSrv4_OK_With_Need_Srv1})
		assert.Nil(t, err)
		_, err = c.Build(reflect.TypeFor[Service1]())
		assert.ErrorIs(t, err, ErrCircularDependency)
		assert.Contains(t, err.Error(),
			"ErrCircularDependency: circular dependency detected at type 'autowire.Service1'")
	})

	t.Run("Circular dependency (self dependent)", func(t *testing.T) {
		// S1 -> S1
		c, err := NewContainer([]any{NewSrv1_Fail_Need_Srv1})
		assert.Nil(t, err)
		_, err = c.Build(reflect.TypeFor[Service1]())
		assert.ErrorIs(t, err, ErrCircularDependency)
		assert.Contains(t, err.Error(),
			"ErrCircularDependency: circular dependency detected at type 'autowire.Service1'")
	})

	t.Run("Requires context.Context, but not provide", func(t *testing.T) {
		c, err := NewContainer([]any{NewSrv1_OK_With_Need_Ctx})
		assert.Nil(t, err)
		_, err = c.Build(reflect.TypeFor[Service1]())
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Contains(t, err.Error(),
			"ErrNotFound: provider not found for type 'context.Context'")
	})
}

func TestContainerBuild_Success(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		c, err := NewContainer([]any{NewSrv1_OK_With_Need_Srv2_Srv3_IntSlice, NewSrv2_OK_With_Need_Srv4_Srv5,
			NewSrv3_OK, NewSrv4_OK, NewSrv5_OK, &struct1_OK, &struct5_OK})
		assert.Nil(t, err)
		s1, err := c.Build(reflect.TypeFor[Service1]())
		assert.Nil(t, err)
		s2, err := c.Get(reflect.TypeFor[Service2]())
		assert.Nil(t, err)
		s3, err := c.Get(reflect.TypeFor[Service3]())
		assert.Nil(t, err)
		assert.Equal(t, []any{s2.Interface().(Service2), s3.Interface().(Service3), struct1_OK.Slice},
			s1.Interface().(Service1).InitArgs())
	})

	t.Run("Success, with context passed", func(t *testing.T) {
		ctx := context.Background()
		c, err := NewContainer([]any{NewSrv1_OK_With_Need_Ctx, NewSrv4_OK, NewSrv5_OK, &struct1_OK})
		assert.Nil(t, err)
		s1, err := c.BuildWithCtx(ctx, reflect.TypeFor[Service1]())
		assert.Nil(t, err)
		assert.Equal(t, []any{ctx}, s1.Interface().(Service1).InitArgs())
		s4, err := c.Build(reflect.TypeFor[Service4]())
		assert.Nil(t, err)
		assert.Equal(t, 0, len(s4.Interface().(Service4).InitArgs()))
		s5, err := c.Build(reflect.TypeFor[Service5]())
		assert.Nil(t, err)
		assert.Equal(t, 0, len(s5.Interface().(Service5).InitArgs()))
	})

	t.Run("Success, with non-shared mode initial", func(t *testing.T) {
		c, err := NewContainer([]any{NewSrv1_OK_With_Need_Srv2_Srv3_IntSlice, NewSrv2_OK_With_Need_Srv4_Srv5,
			NewSrv3_OK, NewSrv4_OK, NewSrv5_OK, &struct1_OK, &struct5_OK}, SetSharedMode(false))
		assert.Nil(t, err)
		s1, err := c.Build(reflect.TypeFor[Service1]())
		assert.Nil(t, err)
		assert.NotNil(t, s1.Interface().(Service1))
		_, err = c.Get(reflect.TypeFor[Service2]())
		assert.ErrorIs(t, err, ErrNotFound)
		s2, err := c.Build(reflect.TypeFor[Service2]())
		assert.Nil(t, err)
		assert.NotNil(t, s2.Interface().(Service2))
	})

	t.Run("Success, with non-shared mode on the fly", func(t *testing.T) {
		c, err := NewContainer([]any{NewSrv1_OK_With_Need_Srv2_Srv3_IntSlice, NewSrv2_OK_With_Need_Srv4_Srv5,
			NewSrv3_OK, NewSrv4_OK, NewSrv5_OK, &struct1_OK, &struct5_OK})
		assert.Nil(t, err)
		s1, err := c.Build(reflect.TypeFor[Service1](), NonSharedMode())
		assert.Nil(t, err)
		assert.NotNil(t, s1.Interface().(Service1))
		_, err = c.Get(reflect.TypeFor[Service2]())
		assert.ErrorIs(t, err, ErrNotFound)
		s2, err := c.Build(reflect.TypeFor[Service2](), NonSharedMode())
		assert.Nil(t, err)
		assert.NotNil(t, s2.Interface().(Service2))
	})

	t.Run("Success, with overwriting type []int", func(t *testing.T) {
		c, err := NewContainer([]any{NewSrv1_OK_With_Need_Srv2_Srv3_IntSlice, NewSrv2_OK, NewSrv3_OK, &struct1_OK})
		assert.Nil(t, err)
		s1, err := c.Build(reflect.TypeFor[Service1](), ProviderOverwrite([]int{1, 1}))
		assert.Nil(t, err)
		s2, err := c.Get(reflect.TypeFor[Service2]())
		assert.Nil(t, err)
		s3, err := c.Get(reflect.TypeFor[Service3]())
		assert.Nil(t, err)
		assert.Equal(t, []any{s2.Interface().(Service2), s3.Interface().(Service3), []int{1, 1}},
			s1.Interface().(Service1).InitArgs())
	})

	t.Run("Success, with updating struct field value", func(t *testing.T) {
		struct1Copy := struct1_OK
		c, err := NewContainer([]any{NewSrv1_OK_With_Need_Srv2_Srv3_IntSlice, NewSrv2_OK, NewSrv3_OK, &struct1Copy})
		assert.Nil(t, err)

		// Update struct field value
		struct1Copy.Slice = []int{100, 200}

		s1, err := c.Build(reflect.TypeFor[Service1]())
		assert.Nil(t, err)
		s2, err := c.Get(reflect.TypeFor[Service2]())
		assert.Nil(t, err)
		s3, err := c.Get(reflect.TypeFor[Service3]())
		assert.Nil(t, err)
		assert.Equal(t, []any{s2.Interface().(Service2), s3.Interface().(Service3), []int{100, 200}},
			s1.Interface().(Service1).InitArgs())
	})
}
