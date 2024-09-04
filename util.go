package autowire

import "reflect"

// typeFor returns the [Type] that represents the type argument T.
func typeFor[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
