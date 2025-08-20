package pointer_utility

import "errors"

func NewPointer[V any](value V) *V {
	return &value
}

func PointerIsNil[P any](pointer *P, msg string) error {
	if pointer == nil {
		return errors.New(msg)
	}
	return nil
}
