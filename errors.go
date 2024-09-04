package autowire

import (
	"errors"
)

var (
	ErrTypeCast           = errors.New("ErrTypeCast")
	ErrNotFound           = errors.New("ErrNotFound")
	ErrProviderInvalid    = errors.New("ErrProviderInvalid")
	ErrProviderDuplicated = errors.New("ErrProviderDuplicated")
	ErrCircularDependency = errors.New("ErrCircularDependency")
)
