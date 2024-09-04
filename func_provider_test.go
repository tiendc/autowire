package autowire

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncProvider_Failure(t *testing.T) {
	t.Run("Func provider returns invalid number of values (3)", func(t *testing.T) {
		_, err := parseProviders(NewSrv1_Fail_With_3_Ret_Values, NewSrv2_OK)
		assert.ErrorIs(t, err, ErrProviderInvalid)
		assert.Contains(t, err.Error(),
			"ErrProviderInvalid: function must return either a <value> or a pair of (<value>, <error>)")
	})

	t.Run("Func provider returns invalid number of values (0)", func(t *testing.T) {
		_, err := parseProviders(NewSrv1_Fail_With_0_Ret_Value, NewSrv2_OK)
		assert.ErrorIs(t, err, ErrProviderInvalid)
		assert.Contains(t, err.Error(),
			"ErrProviderInvalid: function must return either a <value> or a pair of (<value>, <error>)")
	})

	t.Run("Func provider returns invalid type of last value", func(t *testing.T) {
		_, err := parseProviders(NewSrv2_OK, NewSrv1_Fail_With_Non_Err_At_Last, &struct1_OK)
		assert.ErrorIs(t, err, ErrProviderInvalid)
		assert.Contains(t, err.Error(),
			"ErrProviderInvalid: function must return second value of type error")
	})

	t.Run("Func provider has duplicated arg type", func(t *testing.T) {
		_, err := parseProviders(NewSrv2_OK, NewSrv1_Fail_With_Dup_Arg_Type, &struct1_OK)
		assert.ErrorIs(t, err, ErrProviderInvalid)
		assert.Contains(t, err.Error(),
			"ErrProviderInvalid: duplicated function argument type")
	})

	t.Run("Func providers have duplicated returning type", func(t *testing.T) {
		_, err := parseProviders(NewSrv1_OK, NewSrv1_OK_With_Nil_Err)
		assert.ErrorIs(t, err, ErrProviderDuplicated)
		assert.Contains(t, err.Error(),
			"ErrProviderDuplicated: duplicated provider for type 'autowire.Service1'")
	})

	t.Run("Func provider is variadic", func(t *testing.T) {
		_, err := parseProviders(NewSrv1_Fail_With_Variadic, NewSrv1_OK_With_Nil_Err)
		assert.ErrorIs(t, err, ErrProviderInvalid)
		assert.Contains(t, err.Error(),
			"ErrProviderInvalid: variadic function is not allowed")
	})
}

func TestFuncProvider_Success(t *testing.T) {
	t.Run("Simple case", func(t *testing.T) {
		ps1, err := parseProviders(NewSrv1_OK)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(ps1.GetAll()))
		assert.Equal(t, typeFor[Service1](), ps1.GetAll()[0].TargetTypes()[0])
		assert.Equal(t, reflect.TypeOf(NewSrv1_OK), reflect.TypeOf(ps1.GetAll()[0].Source()))
	})
}
