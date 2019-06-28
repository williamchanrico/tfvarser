package tfvarser

import "errors"

var (
	// ErrObjNotSupported this object is not supported
	ErrObjNotSupported = errors.New("Object in this provider is not supported")

	// ErrProviderNotSupported this provider is not supported
	ErrProviderNotSupported = errors.New("Provider is not supported")
)
