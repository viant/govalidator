package govalidator

import (
	"context"
	"fmt"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

type Zeroable interface {
	IsZero() bool
}

func checkRequiredPtr(ctx context.Context, value interface{}) (bool, error) {
	ptr := xunsafe.AsPointer(value)
	if ptr == nil || (*unsafe.Pointer)(ptr) == nil {
		return false, nil
	}
	return true, nil
}

func checkRequiredString(ctx context.Context, value interface{}) (bool, error) {
	text, _ := value.(string)
	return text != "", nil
}

func checkRequiredSlice(ctx context.Context, value interface{}) (bool, error) {
	if value == nil {
		return false, nil
	}
	ptr := xunsafe.AsPointer(value)
	header := (*reflect.SliceHeader)(ptr)
	return header.Len > 0, nil
}

func checkRequiredZeroStruct(ctx context.Context, value interface{}) (bool, error) {
	zeroer, ok := value.(Zeroable)
	if !ok {
		return false, fmt.Errorf("expected: %T, but had: %T", zeroer, value)
	}
	return zeroer.IsZero(), nil
}

func newRequiredCheck(field *Field, check *Check) (IsValid, error) {
	switch field.Kind() {
	case reflect.Ptr:
		return checkRequiredPtr, nil
	case reflect.Struct:
		_, ok := field.Type.MethodByName("Zeroable")
		if !ok {
			return nil, fmt.Errorf("struct does not implemt Zeroable for required check: %v", field.Type.String())
		}
		return checkRequiredZeroStruct, nil
	case reflect.String:
		return checkRequiredString, nil
	case reflect.Int:
		return checkRequiredInt, nil
	case reflect.Slice:
		return checkRequiredSlice, nil
	}
	return nil, fmt.Errorf("required unsupported type: %v %v", field.Name, field.Type.String())
}

func checkRequiredInt(ctx context.Context, value interface{}) (bool, error) {
	intValue, _ := value.(int)
	return intValue != 0, nil
}
