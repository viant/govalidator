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

func checkRequiredNoZeroStruct(ctx context.Context, value interface{}) (bool, error) {
	zeroer, ok := value.(Zeroable)
	if !ok {
		return false, fmt.Errorf("expected: %T, but had: %T", zeroer, value)
	}
	return !zeroer.IsZero(), nil
}

func newRequiredCheck(field *Field, check *Check) (IsValid, error) {
	switch field.Kind() {
	case reflect.Ptr:
		return checkRequiredPtr, nil
	case reflect.Struct:
		_, ok := field.Type.MethodByName("IsZero")
		if !ok {
			return nil, fmt.Errorf("struct does not implement IsZero for required check: %v", field.Type.String())
		}
		return checkRequiredNoZeroStruct, nil
	case reflect.String:
		return checkRequiredString, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return checkRequiredNumeric, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return checkRequiredNumeric, nil
	case reflect.Bool:
		return checkRequiredBool, nil

	case reflect.Float32, reflect.Float64:
		return checkRequiredNumeric, nil
	case reflect.Slice:
		return checkRequiredSlice, nil
	}
	return nil, fmt.Errorf("required unsupported type: %v %v", field.Name, field.Type.String())
}

func checkRequiredNumeric(ctx context.Context, value interface{}) (bool, error) {
	if value == nil {
		return false, nil
	}
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() != 0, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint() != 0, nil
	case reflect.Float32, reflect.Float64:
		return rv.Float() != 0, nil
	}
	return false, fmt.Errorf("unsupported numeric value: %T", value)
}

var boolType = reflect.TypeOf(true)

func checkRequiredBool(ctx context.Context, value interface{}) (bool, error) {
	switch actual := value.(type) {
	case bool:
		return actual, nil
	default:
		ret := reflect.ValueOf(value).Bool()
		return ret, nil
	}
}
