package autowire

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProviderSet_Failure(t *testing.T) {
	t.Run("Empty provider input", func(t *testing.T) {
		_, err := NewProviderSet()
		assert.ErrorIs(t, err, ErrProviderInvalid)
		assert.Contains(t, err.Error(), "ErrProviderInvalid: no provider provided")
	})

	t.Run("Pass nil as provider", func(t *testing.T) {
		_, err := NewProviderSet(NewSrv1_OK, nil)
		assert.ErrorIs(t, err, ErrProviderInvalid)
		assert.Contains(t, err.Error(), "ErrProviderInvalid: provider must not be nil")
	})

	t.Run("Pass invalid provider type (a map)", func(t *testing.T) {
		_, err := NewProviderSet(NewSrv1_OK, map[int]int{})
		assert.ErrorIs(t, err, ErrProviderInvalid)
		assert.Contains(t, err.Error(),
			"ErrProviderInvalid: provider type unsupported, got 'map[int]int'")
	})

	t.Run("On failure, MustNewProviderSet panics", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				assert.False(t, false)
			}
		}()
		ps1 := MustNewProviderSet(NewSrv1_OK, nil)
		assert.Equal(t, 0, len(ps1.GetAll()))
	})

	t.Run("Duplicated provider", func(t *testing.T) {
		ps1, err := NewProviderSet(NewSrv1_OK)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(ps1.GetAll()))
		_, err = NewProviderSet(NewSrv1_OK, ps1.GetAll()[0])
		assert.ErrorIs(t, err, ErrProviderDuplicated)
		assert.Contains(t, err.Error(),
			"ErrProviderDuplicated: duplicated provider for type 'autowire.Service1'")
	})

	t.Run("Duplicated provider with passing a provider set as arg", func(t *testing.T) {
		ps1, err := NewProviderSet(NewSrv1_OK, NewSrv3_OK)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(ps1.GetAll()))
		_, err = NewProviderSet(NewSrv3_OK, ps1)
		assert.ErrorIs(t, err, ErrProviderDuplicated)
	})
}

func TestNewProviderSet_Success(t *testing.T) {
	t.Run("On success, no panic from MustNewProviderSet", func(t *testing.T) {
		defer func() {
			if recover() != nil {
				assert.False(t, false)
			}
		}()
		ps1 := MustNewProviderSet(NewSrv1_OK)
		assert.Equal(t, 1, len(ps1.GetAll()))
	})

	t.Run("Pass a provider as arg", func(t *testing.T) {
		ps1, err := NewProviderSet(NewSrv1_OK)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(ps1.GetAll()))
		ps2, err := NewProviderSet(NewSrv2_OK, ps1.GetAll()[0])
		assert.Nil(t, err)
		assert.Equal(t, 2, len(ps2.GetAll()))
	})

	t.Run("Pass a provider set as arg", func(t *testing.T) {
		ps1, err := NewProviderSet(NewSrv1_OK)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(ps1.GetAll()))
		ps2, err := NewProviderSet(NewSrv2_OK, ps1)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(ps2.GetAll()))
	})

	t.Run("Overwrite a value", func(t *testing.T) {
		ps1, err := NewProviderSet(NewSrv1_OK)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(ps1.GetAll()))
		// Overwrite an existing type
		ps1.Overwrite(newValueProvider(NewSrv1_OK(), reflect.ValueOf(NewSrv1_OK())))
		assert.Equal(t, 1, len(ps1.GetAll()))
		// Overwrite a non-existing type
		ps1.Overwrite(newValueProvider(123, reflect.ValueOf(123)))
		assert.Equal(t, 2, len(ps1.GetAll()))
	})
}
