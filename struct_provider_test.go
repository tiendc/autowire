package autowire

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructProvider_Failure(t *testing.T) {
	t.Run("Pass struct instead of struct pointer", func(t *testing.T) {
		_, err := parseProviders(NewSrv1_OK, struct1_OK)
		assert.ErrorIs(t, err, ErrProviderInvalid)
		assert.Contains(t, err.Error(),
			"ErrProviderInvalid: struct pointer required, got 'autowire.Struct1_OK'")
	})

	t.Run("Pass a nil struct pointer", func(t *testing.T) {
		var structPtr *Struct1_OK
		_, err := parseProviders(NewSrv1_OK, structPtr)
		assert.ErrorIs(t, err, ErrProviderInvalid)
		assert.Contains(t, err.Error(), "ErrProviderInvalid: provider must not be nil")
	})

	t.Run("Struct has duplicated field type (including nested fields)", func(t *testing.T) {
		_, err := parseProviders(NewSrv1_OK, &struct2_Dup_Field_Type)
		assert.ErrorIs(t, err, ErrProviderDuplicated)
		assert.Contains(t, err.Error(), "ErrProviderDuplicated: duplicated provider for type '[]int'")
	})
}

func TestStructProvider_Success(t *testing.T) {
	t.Run("Multiple structs without same type", func(t *testing.T) {
		ps1, err := parseProviders(NewSrv1_OK, &struct1_OK, &struct3_Empty, &struct5_OK)
		assert.Nil(t, err)
		assert.Equal(t, 7, len(ps1.GetAll()))
	})

	t.Run("Struct has duplicated field type but unexported", func(t *testing.T) {
		_, err := parseProviders(NewSrv1_OK, &struct4_OK_Dup_Field_Type_Unexported, &struct3_Empty)
		assert.Nil(t, err)
	})

	t.Run("Struct field has no dependent type", func(t *testing.T) {
		ps1, err := parseProviders(&struct5_OK)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(ps1.GetAll()))
		assert.Equal(t, 0, len(ps1.GetAll()[0].DependentTypes()))
	})

	t.Run("Deep anonymous nested struct", func(t *testing.T) {
		ps1, err := parseProviders(&struct6_OK_Nested_Anonymous)
		assert.Nil(t, err)
		assert.Equal(t, 7, len(ps1.GetAll())) // 5 fields and 2 nested structs themselves
	})

	t.Run("Deep nested struct", func(t *testing.T) {
		ps1, err := parseProviders(&struct7_OK_Nested)
		assert.Nil(t, err)
		assert.Equal(t, 7, len(ps1.GetAll())) // 5 fields and 2 nested structs themselves
		prov1, err := ps1.GetFor(reflect.TypeFor[Nested1]())
		assert.Nil(t, err)
		assert.Equal(t, prov1.Source(), &struct7_OK_Nested)
		prov2, err := ps1.GetFor(reflect.TypeFor[*Nested2]())
		assert.Nil(t, err)
		assert.Equal(t, prov2.Source(), &struct7_OK_Nested)
	})
}
