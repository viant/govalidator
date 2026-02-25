package govalidator

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type boundsCheck struct {
	min float64
	max float64
}

func (b *boundsCheck) minNumeric(ctx context.Context, value interface{}) (bool, error) {
	actual, ok := numericValue(value)
	if !ok {
		return false, nil
	}
	return actual >= b.min, nil
}

func (b *boundsCheck) maxNumeric(ctx context.Context, value interface{}) (bool, error) {
	actual, ok := numericValue(value)
	if !ok {
		return false, nil
	}
	return actual <= b.max, nil
}

func (b *boundsCheck) betweenNumeric(ctx context.Context, value interface{}) (bool, error) {
	actual, ok := numericValue(value)
	if !ok {
		return false, nil
	}
	return actual >= b.min && actual <= b.max, nil
}

func (b *boundsCheck) minLen(ctx context.Context, value interface{}) (bool, error) {
	actual, ok := valueLength(value)
	if !ok {
		return false, nil
	}
	return actual >= int(b.min), nil
}

func (b *boundsCheck) maxLen(ctx context.Context, value interface{}) (bool, error) {
	actual, ok := valueLength(value)
	if !ok {
		return false, nil
	}
	return actual <= int(b.max), nil
}

func (b *boundsCheck) betweenLen(ctx context.Context, value interface{}) (bool, error) {
	actual, ok := valueLength(value)
	if !ok {
		return false, nil
	}
	return actual >= int(b.min) && actual <= int(b.max), nil
}

func NewMin() func(field *Field, check *Check) (IsValid, error) {
	return func(field *Field, check *Check) (IsValid, error) {
		if len(check.Parameters) != 1 {
			return nil, fmt.Errorf("min expects 1 parameter, but had: %d", len(check.Parameters))
		}
		param, err := strconv.ParseFloat(check.Parameters[0], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid min parameter: %w", err)
		}
		ret := &boundsCheck{min: param}
		kind, elemKind := typeKinds(field)

		switch kind {
		case reflect.String:
			if !isWholeNumber(param) {
				return nil, fmt.Errorf("min for string expects integer, but had: %v", check.Parameters[0])
			}
			return ret.minLen, nil
		case reflect.Slice:
			switch elemKind {
			case reflect.String:
				if !isWholeNumber(param) {
					return nil, fmt.Errorf("min for string expects integer, but had: %v", check.Parameters[0])
				}
				return ret.minLen, nil
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Float32, reflect.Float64:
				return ret.minNumeric, nil
			case reflect.Slice, reflect.Struct, reflect.Ptr:
				if !isWholeNumber(param) {
					return nil, fmt.Errorf("min for slice expects integer, but had: %v", check.Parameters[0])
				}
				return ret.minLen, nil
			}
			if !isWholeNumber(param) {
				return nil, fmt.Errorf("min for slice expects integer, but had: %v", check.Parameters[0])
			}
			return ret.minLen, nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			return ret.minNumeric, nil
		}
		return nil, fmt.Errorf("unsupported min type: %s", field.Type.String())
	}
}

func NewMax() func(field *Field, check *Check) (IsValid, error) {
	return func(field *Field, check *Check) (IsValid, error) {
		if len(check.Parameters) != 1 {
			return nil, fmt.Errorf("max expects 1 parameter, but had: %d", len(check.Parameters))
		}
		param, err := strconv.ParseFloat(check.Parameters[0], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid max parameter: %w", err)
		}
		ret := &boundsCheck{max: param}
		kind, elemKind := typeKinds(field)

		switch kind {
		case reflect.String:
			if !isWholeNumber(param) {
				return nil, fmt.Errorf("max for string expects integer, but had: %v", check.Parameters[0])
			}
			return ret.maxLen, nil
		case reflect.Slice:
			switch elemKind {
			case reflect.String:
				if !isWholeNumber(param) {
					return nil, fmt.Errorf("max for string expects integer, but had: %v", check.Parameters[0])
				}
				return ret.maxLen, nil
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Float32, reflect.Float64:
				return ret.maxNumeric, nil
			case reflect.Slice, reflect.Struct, reflect.Ptr:
				if !isWholeNumber(param) {
					return nil, fmt.Errorf("max for slice expects integer, but had: %v", check.Parameters[0])
				}
				return ret.maxLen, nil
			}
			if !isWholeNumber(param) {
				return nil, fmt.Errorf("max for slice expects integer, but had: %v", check.Parameters[0])
			}
			return ret.maxLen, nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			return ret.maxNumeric, nil
		}
		return nil, fmt.Errorf("unsupported max type: %s", field.Type.String())
	}
}

func NewBetween() func(field *Field, check *Check) (IsValid, error) {
	return func(field *Field, check *Check) (IsValid, error) {
		if len(check.Parameters) != 2 {
			return nil, fmt.Errorf("between expects 2 parameters, but had: %d", len(check.Parameters))
		}
		min, err := strconv.ParseFloat(check.Parameters[0], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid between min parameter: %w", err)
		}
		max, err := strconv.ParseFloat(check.Parameters[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid between max parameter: %w", err)
		}
		ret := &boundsCheck{min: min, max: max}
		kind, elemKind := typeKinds(field)

		switch kind {
		case reflect.String:
			if !isWholeNumber(min) || !isWholeNumber(max) {
				return nil, fmt.Errorf("between for string expects integers, but had: %v,%v", check.Parameters[0], check.Parameters[1])
			}
			return ret.betweenLen, nil
		case reflect.Slice:
			switch elemKind {
			case reflect.String:
				if !isWholeNumber(min) || !isWholeNumber(max) {
					return nil, fmt.Errorf("between for string expects integers, but had: %v,%v", check.Parameters[0], check.Parameters[1])
				}
				return ret.betweenLen, nil
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Float32, reflect.Float64:
				return ret.betweenNumeric, nil
			case reflect.Slice, reflect.Struct, reflect.Ptr:
				if !isWholeNumber(min) || !isWholeNumber(max) {
					return nil, fmt.Errorf("between for slice expects integers, but had: %v,%v", check.Parameters[0], check.Parameters[1])
				}
				return ret.betweenLen, nil
			}
			if !isWholeNumber(min) || !isWholeNumber(max) {
				return nil, fmt.Errorf("between for slice expects integers, but had: %v,%v", check.Parameters[0], check.Parameters[1])
			}
			return ret.betweenLen, nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			return ret.betweenNumeric, nil
		}
		return nil, fmt.Errorf("unsupported between type: %s", field.Type.String())
	}
}

type stringCheck struct {
	arg string
}

func (s *stringCheck) contains(ctx context.Context, value interface{}) (bool, error) {
	actual, ok := asStringValue(value)
	if !ok {
		return false, nil
	}
	return strings.Contains(actual, s.arg), nil
}

func (s *stringCheck) notContains(ctx context.Context, value interface{}) (bool, error) {
	actual, ok := asStringValue(value)
	if !ok {
		return false, nil
	}
	return !strings.Contains(actual, s.arg), nil
}

func (s *stringCheck) startsWith(ctx context.Context, value interface{}) (bool, error) {
	actual, ok := asStringValue(value)
	if !ok {
		return false, nil
	}
	return strings.HasPrefix(actual, s.arg), nil
}

func (s *stringCheck) endsWith(ctx context.Context, value interface{}) (bool, error) {
	actual, ok := asStringValue(value)
	if !ok {
		return false, nil
	}
	return strings.HasSuffix(actual, s.arg), nil
}

func NewContains() func(field *Field, check *Check) (IsValid, error) {
	return newStringCheck(func(s *stringCheck) IsValid {
		return s.contains
	}, "contains")
}

func NewNotContains() func(field *Field, check *Check) (IsValid, error) {
	return newStringCheck(func(s *stringCheck) IsValid {
		return s.notContains
	}, "notcontains")
}

func NewStartsWith() func(field *Field, check *Check) (IsValid, error) {
	return newStringCheck(func(s *stringCheck) IsValid {
		return s.startsWith
	}, "startswith")
}

func NewEndsWith() func(field *Field, check *Check) (IsValid, error) {
	return newStringCheck(func(s *stringCheck) IsValid {
		return s.endsWith
	}, "endswith")
}

func newStringCheck(factory func(*stringCheck) IsValid, name string) func(field *Field, check *Check) (IsValid, error) {
	return func(field *Field, check *Check) (IsValid, error) {
		if len(check.Parameters) != 1 {
			return nil, fmt.Errorf("%s expects 1 parameter, but had: %d", name, len(check.Parameters))
		}
		kind, elemKind := typeKinds(field)
		switch kind {
		case reflect.String:
		case reflect.Slice:
			if elemKind != reflect.String {
				return nil, fmt.Errorf("unsupported %s type: %s", name, field.Type.String())
			}
		default:
			return nil, fmt.Errorf("unsupported %s type: %s", name, field.Type.String())
		}
		return factory(&stringCheck{arg: check.Parameters[0]}), nil
	}
}

type fieldCheck struct {
	other string
}

func (f *fieldCheck) eqField(ctx context.Context, value interface{}) (bool, error) {
	otherValue, err := lookupOtherField(ctx, f.other)
	if err != nil {
		return false, err
	}
	return equalValues(value, otherValue), nil
}

func (f *fieldCheck) neField(ctx context.Context, value interface{}) (bool, error) {
	otherValue, err := lookupOtherField(ctx, f.other)
	if err != nil {
		return false, err
	}
	return !equalValues(value, otherValue), nil
}

func (f *fieldCheck) gtField(ctx context.Context, value interface{}) (bool, error) {
	otherValue, err := lookupOtherField(ctx, f.other)
	if err != nil {
		return false, err
	}
	compare, err := compareValues(value, otherValue)
	if err != nil {
		return false, err
	}
	return compare > 0, nil
}

func NewEqField() func(field *Field, check *Check) (IsValid, error) {
	return newCrossFieldCheck(func(c *fieldCheck) IsValid {
		return c.eqField
	}, "eqfield")
}

func NewNeField() func(field *Field, check *Check) (IsValid, error) {
	return newCrossFieldCheck(func(c *fieldCheck) IsValid {
		return c.neField
	}, "nefield")
}

func NewGtField() func(field *Field, check *Check) (IsValid, error) {
	return newCrossFieldCheck(func(c *fieldCheck) IsValid {
		return c.gtField
	}, "gtfield")
}

func newCrossFieldCheck(factory func(*fieldCheck) IsValid, name string) func(field *Field, check *Check) (IsValid, error) {
	return func(field *Field, check *Check) (IsValid, error) {
		if len(check.Parameters) != 1 {
			return nil, fmt.Errorf("%s expects 1 parameter, but had: %d", name, len(check.Parameters))
		}
		return factory(&fieldCheck{other: check.Parameters[0]}), nil
	}
}

type conditionalRequiredCheck struct {
	otherField   string
	expected     string
	otherFields  []string
	checkCurrent func(current interface{}, condition bool) (bool, error)
}

func (c *conditionalRequiredCheck) requiredIf(ctx context.Context, value interface{}) (bool, error) {
	otherValue, err := lookupOtherField(ctx, c.otherField)
	if err != nil {
		return false, err
	}
	condition := equalsParamValue(otherValue, c.expected)
	return c.checkCurrent(value, condition)
}

func (c *conditionalRequiredCheck) requiredUnless(ctx context.Context, value interface{}) (bool, error) {
	otherValue, err := lookupOtherField(ctx, c.otherField)
	if err != nil {
		return false, err
	}
	condition := !equalsParamValue(otherValue, c.expected)
	return c.checkCurrent(value, condition)
}

func (c *conditionalRequiredCheck) requiredWith(ctx context.Context, value interface{}) (bool, error) {
	condition := false
	for _, other := range c.otherFields {
		otherValue, err := lookupOtherField(ctx, other)
		if err != nil {
			return false, err
		}
		if !isEmpty(otherValue) {
			condition = true
			break
		}
	}
	return c.checkCurrent(value, condition)
}

func (c *conditionalRequiredCheck) requiredWithout(ctx context.Context, value interface{}) (bool, error) {
	condition := false
	for _, other := range c.otherFields {
		otherValue, err := lookupOtherField(ctx, other)
		if err != nil {
			return false, err
		}
		if isEmpty(otherValue) {
			condition = true
			break
		}
	}
	return c.checkCurrent(value, condition)
}

func NewRequiredIf() func(field *Field, check *Check) (IsValid, error) {
	return func(field *Field, check *Check) (IsValid, error) {
		if len(check.Parameters) != 2 {
			return nil, fmt.Errorf("required_if expects 2 parameters, but had: %d", len(check.Parameters))
		}
		c := &conditionalRequiredCheck{
			otherField: check.Parameters[0],
			expected:   check.Parameters[1],
			checkCurrent: func(current interface{}, condition bool) (bool, error) {
				if !condition {
					return true, nil
				}
				return !isEmpty(current), nil
			},
		}
		return c.requiredIf, nil
	}
}

func NewRequiredUnless() func(field *Field, check *Check) (IsValid, error) {
	return func(field *Field, check *Check) (IsValid, error) {
		if len(check.Parameters) != 2 {
			return nil, fmt.Errorf("required_unless expects 2 parameters, but had: %d", len(check.Parameters))
		}
		c := &conditionalRequiredCheck{
			otherField: check.Parameters[0],
			expected:   check.Parameters[1],
			checkCurrent: func(current interface{}, condition bool) (bool, error) {
				if !condition {
					return true, nil
				}
				return !isEmpty(current), nil
			},
		}
		return c.requiredUnless, nil
	}
}

func NewRequiredWith() func(field *Field, check *Check) (IsValid, error) {
	return func(field *Field, check *Check) (IsValid, error) {
		if len(check.Parameters) == 0 {
			return nil, fmt.Errorf("required_with expects at least 1 parameter")
		}
		c := &conditionalRequiredCheck{
			otherFields: check.Parameters,
			checkCurrent: func(current interface{}, condition bool) (bool, error) {
				if !condition {
					return true, nil
				}
				return !isEmpty(current), nil
			},
		}
		return c.requiredWith, nil
	}
}

func NewRequiredWithout() func(field *Field, check *Check) (IsValid, error) {
	return func(field *Field, check *Check) (IsValid, error) {
		if len(check.Parameters) == 0 {
			return nil, fmt.Errorf("required_without expects at least 1 parameter")
		}
		c := &conditionalRequiredCheck{
			otherFields: check.Parameters,
			checkCurrent: func(current interface{}, condition bool) (bool, error) {
				if !condition {
					return true, nil
				}
				return !isEmpty(current), nil
			},
		}
		return c.requiredWithout, nil
	}
}

func NewPast() func(field *Field, check *Check) (IsValid, error) {
	return newTimeRelativeCheck("past", func(actual, now time.Time) bool {
		return actual.Before(now)
	})
}

func NewFuture() func(field *Field, check *Check) (IsValid, error) {
	return newTimeRelativeCheck("future", func(actual, now time.Time) bool {
		return actual.After(now)
	})
}

func typeKinds(field *Field) (reflect.Kind, reflect.Kind) {
	switch field.Kind() {
	case reflect.Ptr:
		elem := field.Elem()
		if elem.Kind() == reflect.Slice {
			return reflect.Slice, elem.Elem().Kind()
		}
		return elem.Kind(), reflect.Invalid
	case reflect.Slice:
		return reflect.Slice, field.Elem().Kind()
	default:
		return field.Kind(), reflect.Invalid
	}
}

func newTimeRelativeCheck(name string, predicate func(actual, now time.Time) bool) func(field *Field, check *Check) (IsValid, error) {
	return func(field *Field, check *Check) (IsValid, error) {
		kind, elemKind := typeKinds(field)
		supported := false
		switch kind {
		case reflect.String, reflect.Struct:
			supported = true
		case reflect.Slice:
			supported = elemKind == reflect.String
		}
		if !supported {
			return nil, fmt.Errorf("unsupported %s type: %s", name, field.Type.String())
		}
		return func(ctx context.Context, value interface{}) (bool, error) {
			actual, ok := toTimeValue(value)
			if !ok {
				return false, nil
			}
			return predicate(actual, time.Now()), nil
		}, nil
	}
}

func isWholeNumber(value float64) bool {
	return math.Trunc(value) == value
}

func numericValue(value interface{}) (float64, bool) {
	if value == nil {
		return 0, false
	}
	rv := reflect.ValueOf(value)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return 0, false
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(rv.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(rv.Uint()), true
	case reflect.Float32, reflect.Float64:
		return rv.Float(), true
	default:
		return 0, false
	}
}

func valueLength(value interface{}) (int, bool) {
	if value == nil {
		return 0, false
	}
	rv := reflect.ValueOf(value)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return 0, false
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.String:
		return utf8.RuneCountInString(rv.String()), true
	case reflect.Slice:
		return rv.Len(), true
	default:
		return 0, false
	}
}

func asStringValue(value interface{}) (string, bool) {
	if value == nil {
		return "", false
	}
	rv := reflect.ValueOf(value)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return "", false
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.String {
		return "", false
	}
	return rv.String(), true
}

func toTimeValue(value interface{}) (time.Time, bool) {
	if value == nil {
		return time.Time{}, false
	}
	rv := reflect.ValueOf(value)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return time.Time{}, false
		}
		rv = rv.Elem()
	}
	if rv.Type() == reflect.TypeOf(time.Time{}) {
		return rv.Interface().(time.Time), true
	}
	if rv.Kind() != reflect.String {
		return time.Time{}, false
	}
	actual := rv.String()
	if actual == "" {
		return time.Time{}, false
	}
	formats := []string{time.RFC3339Nano, time.RFC3339, "2006-01-02"}
	for _, layout := range formats {
		if parsed, err := time.Parse(layout, actual); err == nil {
			return parsed, true
		}
	}
	return time.Time{}, false
}

func lookupOtherField(ctx context.Context, field string) (interface{}, error) {
	session, ok := ctx.Value(SessionKey).(*Session)
	if !ok || session == nil {
		return nil, fmt.Errorf("validation session was not available")
	}
	parent := session.ParentValue
	if parent == nil {
		return nil, fmt.Errorf("validation parent value was not available")
	}
	rv := reflect.ValueOf(parent)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil, nil
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("eqfield family expects struct parent, but had: %T", parent)
	}
	fieldValue := rv.FieldByName(field)
	if !fieldValue.IsValid() {
		return nil, fmt.Errorf("field %q was not found on %T", field, parent)
	}
	if !fieldValue.CanInterface() {
		return nil, fmt.Errorf("field %q is not accessible on %T", field, parent)
	}
	return fieldValue.Interface(), nil
}

func derefReflectValue(value interface{}) (reflect.Value, bool) {
	if value == nil {
		return reflect.Value{}, true
	}
	rv := reflect.ValueOf(value)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return reflect.Value{}, true
		}
		rv = rv.Elem()
	}
	return rv, false
}

func equalValues(left, right interface{}) bool {
	lv, lnil := derefReflectValue(left)
	rv, rnil := derefReflectValue(right)
	if lnil || rnil {
		return lnil && rnil
	}
	if ln, ok := numericValue(left); ok {
		if rn, ok := numericValue(right); ok {
			return ln == rn
		}
	}
	return reflect.DeepEqual(lv.Interface(), rv.Interface())
}

func compareValues(left, right interface{}) (int, error) {
	lv, lnil := derefReflectValue(left)
	rv, rnil := derefReflectValue(right)
	if lnil || rnil {
		return 0, fmt.Errorf("cannot compare nil value")
	}
	if ln, ok := numericValue(left); ok {
		if rn, ok := numericValue(right); ok {
			if ln > rn {
				return 1, nil
			}
			if ln < rn {
				return -1, nil
			}
			return 0, nil
		}
	}
	if lv.Kind() == reflect.String && rv.Kind() == reflect.String {
		ls, rs := lv.String(), rv.String()
		if ls > rs {
			return 1, nil
		}
		if ls < rs {
			return -1, nil
		}
		return 0, nil
	}
	if lv.Type() == reflect.TypeOf(time.Time{}) && rv.Type() == reflect.TypeOf(time.Time{}) {
		lt := lv.Interface().(time.Time)
		rt := rv.Interface().(time.Time)
		if lt.After(rt) {
			return 1, nil
		}
		if lt.Before(rt) {
			return -1, nil
		}
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported compare types: %s and %s", lv.Type(), rv.Type())
}

func equalsParamValue(value interface{}, expected string) bool {
	actual, isNil := derefReflectValue(value)
	if isNil {
		return strings.EqualFold(expected, "nil") || expected == ""
	}
	switch actual.Kind() {
	case reflect.String:
		return actual.String() == expected
	case reflect.Bool:
		parsed, err := strconv.ParseBool(expected)
		if err != nil {
			return false
		}
		return actual.Bool() == parsed
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(expected, 10, 64)
		if err != nil {
			return false
		}
		return actual.Int() == parsed
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsed, err := strconv.ParseUint(expected, 10, 64)
		if err != nil {
			return false
		}
		return actual.Uint() == parsed
	case reflect.Float32, reflect.Float64:
		parsed, err := strconv.ParseFloat(expected, 64)
		if err != nil {
			return false
		}
		return actual.Float() == parsed
	default:
		return fmt.Sprintf("%v", actual.Interface()) == expected
	}
}
