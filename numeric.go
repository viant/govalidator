package govalidator

import (
	"context"
	"fmt"
	"github.com/viant/xunsafe"
	"reflect"
	"strconv"
	"unicode/utf8"
)

type Numeric struct{ param int }

func (n *Numeric) intGt(ctx context.Context, value interface{}) (bool, error) {
	if isBaseTypeNil(value) {
		return false, nil
	}
	return intValue(value) > n.param, nil
}

func (n *Numeric) floatGt(ctx context.Context, value interface{}) (bool, error) {
	if isBaseTypeNil(value) {
		return false, nil
	}
	return floatValue(value) > float64(n.param), nil
}

func (n *Numeric) stringGt(ctx context.Context, value interface{}) (bool, error) {
	if isBaseTypeNil(value) {
		return false, nil
	}
	return int(utf8.RuneCountInString(stringValue(value))) > n.param, nil
}
func (n *Numeric) timeGt(ctx context.Context, value interface{}) (bool, error) {
	return int(utf8.RuneCountInString(stringValue(value))) > n.param, nil
}

func (n *Numeric) sliceGt(ctx context.Context, value interface{}) (bool, error) {
	return sliceLen(value) > n.param, nil
}

func (n *Numeric) intLt(ctx context.Context, value interface{}) (bool, error) {
	if isBaseTypeNil(value) {
		return false, nil
	}
	return intValue(value) < n.param, nil
}

func (n *Numeric) floatLt(ctx context.Context, value interface{}) (bool, error) {
	if isBaseTypeNil(value) {
		return false, nil
	}
	return floatValue(value) < float64(n.param), nil
}

func (n *Numeric) stringLt(ctx context.Context, value interface{}) (bool, error) {
	if isBaseTypeNil(value) {
		return false, nil
	}
	return int(utf8.RuneCountInString(stringValue(value))) < n.param, nil
}
func (n *Numeric) timeLt(ctx context.Context, value interface{}) (bool, error) {
	return int(utf8.RuneCountInString(stringValue(value))) < n.param, nil
}

func (n *Numeric) sliceLt(ctx context.Context, value interface{}) (bool, error) {
	return sliceLen(value) < n.param, nil
}

func (n *Numeric) intGte(ctx context.Context, value interface{}) (bool, error) {
	if isBaseTypeNil(value) {
		return false, nil
	}
	return intValue(value) >= n.param, nil
}

func (n *Numeric) floatGte(ctx context.Context, value interface{}) (bool, error) {
	if isBaseTypeNil(value) {
		return false, nil
	}
	return floatValue(value) >= float64(n.param), nil
}

func (n *Numeric) stringGte(ctx context.Context, value interface{}) (bool, error) {
	if isBaseTypeNil(value) {
		return false, nil
	}
	return int(utf8.RuneCountInString(stringValue(value))) >= n.param, nil
}

func (n *Numeric) timeGte(ctx context.Context, value interface{}) (bool, error) {
	return int(utf8.RuneCountInString(stringValue(value))) >= n.param, nil
}

func (n *Numeric) sliceGte(ctx context.Context, value interface{}) (bool, error) {
	return sliceLen(value) >= n.param, nil
}

func (n *Numeric) intLte(ctx context.Context, value interface{}) (bool, error) {
	if isBaseTypeNil(value) {
		return false, nil
	}
	return intValue(value) <= n.param, nil
}

func (n *Numeric) floatLte(ctx context.Context, value interface{}) (bool, error) {
	if isBaseTypeNil(value) {
		return false, nil
	}
	return floatValue(value) <= float64(n.param), nil
}

func (n *Numeric) stringLte(ctx context.Context, value interface{}) (bool, error) {
	if isBaseTypeNil(value) {
		return false, nil
	}
	return int(utf8.RuneCountInString(stringValue(value))) <= n.param, nil
}

func (n *Numeric) timeLte(ctx context.Context, value interface{}) (bool, error) {
	return int(utf8.RuneCountInString(stringValue(value))) <= n.param, nil
}

func (n *Numeric) sliceLte(ctx context.Context, value interface{}) (bool, error) {
	return sliceLen(value) <= n.param, nil
}

func sliceLen(value interface{}) int {
	header := (*reflect.SliceHeader)(xunsafe.AsPointer(value))
	return header.Len
}
func isBaseTypeNil(value interface{}) bool {
	if value == nil {
		return true
	}
	switch actual := value.(type) {
	case *string:
		return (*string)(xunsafe.AsPointer(actual)) == nil
	case *int, *int64, *uint, *uint64:
		return (*int)(xunsafe.AsPointer(actual)) == nil
	case *uint32, *int32:
		return (*int32)(xunsafe.AsPointer(actual)) == nil
	case *int16, *uint16:
		return (*int16)(xunsafe.AsPointer(actual)) == nil
	case *int8, *uint8:
		return (*int8)(xunsafe.AsPointer(actual)) == nil
	}
	return false
}
func intValue(value interface{}) int {
	switch actual := value.(type) {
	case int:
		return actual
	case int64:
		return int(actual)
	case uint:
		return int(actual)
	case uint64:
		return int(actual)
	case int32:
		return int(actual)
	case uint32:
		return int(actual)
	case int16:
		return int(actual)
	case uint16:
		return int(actual)
	case int8:
		return int(actual)
	case uint8:
		return int(actual)
	case *int, *int64, *uint, *uint64:
		return *(*int)(xunsafe.AsPointer(actual))
	case *uint32, *int32:
		return int(*(*int32)(xunsafe.AsPointer(actual)))
	case *int16, *uint16:
		return int(*(*int16)(xunsafe.AsPointer(actual)))
	case *int8, *uint8:
		return int(*(*int8)(xunsafe.AsPointer(actual)))
	}
	return 0
}

func stringValue(value interface{}) string {
	switch actual := value.(type) {
	case string:
		return actual
	case *string:
		if actual == nil {
			return ""
		}
		return *actual
	}
	return ""
}

func floatValue(value interface{}) float64 {
	switch actual := value.(type) {
	case float64:
		return float64(actual)
	case float32:
		return float64(actual)

	case *float64:
		if actual == nil {
			return 0
		}
		return float64(*actual)
	case *float32:
		if actual == nil {
			return 0
		}
		return float64(*actual)
	}
	return 0
}

// NewGt creates greater than validation check
func NewGt() func(field *Field, check *Check) (IsValid, error) {
	return func(field *Field, check *Check) (IsValid, error) {
		param, err := strconv.Atoi(check.Parameters[0])
		if err != nil {
			return nil, fmt.Errorf("invalid parameter: %v", err)
		}
		ret := &Numeric{param: param}

		switch field.Kind() {
		case reflect.String:
			return ret.stringGt, nil
		case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Int32, reflect.Uint32, reflect.Int16, reflect.Uint16, reflect.Int8, reflect.Uint8:
			return ret.intGt, nil
		case reflect.Float64, reflect.Float32:
			return ret.floatGt, nil
		case reflect.Slice:
			switch field.Elem().Kind() {
			case reflect.String:
				return ret.stringGt, nil
			case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Int32, reflect.Uint32, reflect.Int16, reflect.Uint16, reflect.Int8, reflect.Uint8:
				return ret.intGt, nil
			case reflect.Float64, reflect.Float32:
				return ret.floatGt, nil
			case reflect.Slice:
				return ret.sliceGt, nil
			}
		case reflect.Ptr:
			switch field.Elem().Kind() {
			case reflect.String:
				return ret.stringGt, nil
			case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Int32, reflect.Uint32, reflect.Int16, reflect.Uint16, reflect.Int8, reflect.Uint8:
				return ret.intGt, nil
			case reflect.Float64, reflect.Float32:
				return ret.floatGt, nil
			case reflect.Slice:
				switch field.Elem().Elem().Kind() {
				case reflect.String:
					return ret.stringGt, nil
				case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Int32, reflect.Uint32, reflect.Int16, reflect.Uint16, reflect.Int8, reflect.Uint8:
					return ret.intGt, nil
				case reflect.Float64, reflect.Float32:
					return ret.floatGt, nil
				case reflect.Slice:
					return ret.sliceGt, nil
				}
			}
		}
		return nil, fmt.Errorf("unsupported ge type: %s", field.Type.String())
	}
}

// NewLt creates less than validation check
func NewLt() func(field *Field, check *Check) (IsValid, error) {
	return func(field *Field, check *Check) (IsValid, error) {
		param, err := strconv.Atoi(check.Parameters[0])
		if err != nil {
			return nil, fmt.Errorf("invalid parameter: %v", err)
		}
		ret := &Numeric{param: param}

		switch field.Kind() {
		case reflect.String:
			return ret.stringLt, nil
		case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Int32, reflect.Uint32, reflect.Int16, reflect.Uint16, reflect.Int8, reflect.Uint8:
			return ret.intLt, nil
		case reflect.Float64, reflect.Float32:
			return ret.floatLt, nil
		case reflect.Slice:
			switch field.Elem().Kind() {
			case reflect.String:
				return ret.stringLt, nil
			case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Int32, reflect.Uint32, reflect.Int16, reflect.Uint16, reflect.Int8, reflect.Uint8:
				return ret.intLt, nil
			case reflect.Float64, reflect.Float32:
				return ret.floatLt, nil
			case reflect.Slice:
				return ret.sliceLt, nil
			}
		case reflect.Ptr:
			switch field.Elem().Kind() {
			case reflect.String:
				return ret.stringLt, nil
			case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Int32, reflect.Uint32, reflect.Int16, reflect.Uint16, reflect.Int8, reflect.Uint8:
				return ret.intLt, nil
			case reflect.Float64, reflect.Float32:
				return ret.floatLt, nil
			case reflect.Slice:
				switch field.Elem().Elem().Kind() {
				case reflect.String:
					return ret.stringLt, nil
				case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Int32, reflect.Uint32, reflect.Int16, reflect.Uint16, reflect.Int8, reflect.Uint8:
					return ret.intLt, nil
				case reflect.Float64, reflect.Float32:
					return ret.floatLt, nil
				case reflect.Slice:
					return ret.sliceLt, nil
				}
			}
		}
		return nil, fmt.Errorf("unsupported ge type: %s", field.Type.String())
	}
}

// NewGte creates greater or equal than validation check
func NewGte() func(field *Field, check *Check) (IsValid, error) {
	return func(field *Field, check *Check) (IsValid, error) {
		param, err := strconv.Atoi(check.Parameters[0])
		if err != nil {
			return nil, fmt.Errorf("invalid parameter: %v", err)
		}
		ret := &Numeric{param: param}

		switch field.Kind() {
		case reflect.String:
			return ret.stringGte, nil
		case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Int32, reflect.Uint32, reflect.Int16, reflect.Uint16, reflect.Int8, reflect.Uint8:
			return ret.intGte, nil
		case reflect.Float64, reflect.Float32:
			return ret.floatGte, nil
		case reflect.Slice:
			switch field.Elem().Kind() {
			case reflect.String:
				return ret.stringGte, nil
			case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Int32, reflect.Uint32, reflect.Int16, reflect.Uint16, reflect.Int8, reflect.Uint8:
				return ret.intGte, nil
			case reflect.Float64, reflect.Float32:
				return ret.floatGte, nil
			case reflect.Slice:
				return ret.sliceGte, nil
			}
		case reflect.Ptr:
			switch field.Elem().Kind() {
			case reflect.String:
				return ret.stringGte, nil
			case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Int32, reflect.Uint32, reflect.Int16, reflect.Uint16, reflect.Int8, reflect.Uint8:
				return ret.intGte, nil
			case reflect.Float64, reflect.Float32:
				return ret.floatGte, nil
			case reflect.Slice:
				switch field.Elem().Elem().Kind() {
				case reflect.String:
					return ret.stringGte, nil
				case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Int32, reflect.Uint32, reflect.Int16, reflect.Uint16, reflect.Int8, reflect.Uint8:
					return ret.intGte, nil
				case reflect.Float64, reflect.Float32:
					return ret.floatGte, nil
				case reflect.Slice:
					return ret.sliceGte, nil
				}
			}
		}
		return nil, fmt.Errorf("unsupported ge type: %s", field.Type.String())
	}
}

// NewLte creates less or equal than validation check
func NewLte() func(field *Field, check *Check) (IsValid, error) {
	return func(field *Field, check *Check) (IsValid, error) {
		param, err := strconv.Atoi(check.Parameters[0])
		if err != nil {
			return nil, fmt.Errorf("invalid parameter: %v", err)
		}
		ret := &Numeric{param: param}

		switch field.Kind() {
		case reflect.String:
			return ret.stringLte, nil
		case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Int32, reflect.Uint32, reflect.Int16, reflect.Uint16, reflect.Int8, reflect.Uint8:
			return ret.intLte, nil
		case reflect.Float64, reflect.Float32:
			return ret.floatLte, nil
		case reflect.Slice:
			switch field.Elem().Kind() {
			case reflect.String:
				return ret.stringLte, nil
			case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Int32, reflect.Uint32, reflect.Int16, reflect.Uint16, reflect.Int8, reflect.Uint8:
				return ret.intLte, nil
			case reflect.Float64, reflect.Float32:
				return ret.floatLte, nil
			case reflect.Slice:
				return ret.sliceLte, nil
			}
		case reflect.Ptr:
			switch field.Elem().Kind() {
			case reflect.String:
				return ret.stringLte, nil
			case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Int32, reflect.Uint32, reflect.Int16, reflect.Uint16, reflect.Int8, reflect.Uint8:
				return ret.intLte, nil
			case reflect.Float64, reflect.Float32:
				return ret.floatLte, nil
			case reflect.Slice:
				switch field.Elem().Elem().Kind() {
				case reflect.String:
					return ret.stringLte, nil
				case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Int32, reflect.Uint32, reflect.Int16, reflect.Uint16, reflect.Int8, reflect.Uint8:
					return ret.intLte, nil
				case reflect.Float64, reflect.Float32:
					return ret.floatLte, nil
				case reflect.Slice:
					return ret.sliceLte, nil
				}
			}
		}
		return nil, fmt.Errorf("unsupported ge type: %s", field.Type.String())
	}
}
