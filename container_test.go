package autowire

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContainer_Failure(t *testing.T) {
	t.Run("No provider provided", func(t *testing.T) {
		_, err := NewContainer(nil)
		assert.ErrorIs(t, err, ErrProviderInvalid)
		assert.Contains(t, err.Error(), "ErrProviderInvalid: no provider provided")
	})

	t.Run("No provider provided, MustNewContainer panics", func(t *testing.T) {
		defer func() {
			err := recover().(error)
			assert.ErrorIs(t, err, ErrProviderInvalid)
			assert.Contains(t, err.Error(), "ErrProviderInvalid: no provider provided")
		}()
		_ = MustNewContainer(nil)
	})
}

func TestNewContainer_Success(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		c, err := NewContainer([]any{NewSrv1_OK})
		assert.Nil(t, err)
		assert.Equal(t, true, c.SharedMode())
		assert.Equal(t, reflect.TypeOf(NewSrv1_OK), reflect.TypeOf(c.ProviderSet().GetAll()[0].Source()))
	})

	t.Run("Success MustNewContainer", func(t *testing.T) {
		defer func() {
			if recover() != nil {
				assert.False(t, false)
			}
		}()
		c := MustNewContainer([]any{NewSrv1_OK})
		assert.Equal(t, true, c.SharedMode())
		assert.Equal(t, reflect.TypeOf(NewSrv1_OK), reflect.TypeOf(c.ProviderSet().GetAll()[0].Source()))
	})
}

func TestContainerGet(t *testing.T) {
	t.Run("Not found", func(t *testing.T) {
		c, err := NewContainer([]any{NewSrv1_OK})
		assert.Nil(t, err)
		_, err = c.Get(typeFor[Service1]())
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Contains(t, err.Error(), "ErrNotFound: object not found for type 'autowire.Service1'")
	})

	t.Run("Success", func(t *testing.T) {
		c, err := NewContainer([]any{NewSrv1_OK})
		assert.Nil(t, err)
		v1, err := c.Build(typeFor[Service1]())
		assert.Nil(t, err)
		v2, err := c.Get(typeFor[Service1]())
		assert.Nil(t, err)
		assert.Equal(t, v1, v2)
		_, err = c.Get(typeFor[Service2]())
		assert.ErrorIs(t, err, ErrNotFound)
	})
}
