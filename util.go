package memberclicks

import (
	"errors"
	"reflect"
)

// Errors related to copying values between interfaces
var (
	ErrInvalidDstVal     = errors.New("invalid destination value")
	ErrCannotAssignValue = errors.New("cannot assign to destination value")
)

// copy copies one interface into the other doing type checking to make sure
// it's safe. If it cannot be copied, an error is returned.
func copy(srcVal interface{}, dstVal interface{}) error {

	curEl := reflect.ValueOf(srcVal)
	if curEl.Kind() == reflect.Ptr {
		curEl = curEl.Elem()
	}
	dstEl := reflect.ValueOf(dstVal)
	if dstEl.Kind() == reflect.Ptr {
		dstEl = dstEl.Elem()
	}
	if !dstEl.CanSet() {
		return ErrInvalidDstVal
	}
	if !curEl.Type().AssignableTo(dstEl.Type()) {
		return ErrCannotAssignValue
	}
	dstEl.Set(curEl)
	return nil
}
