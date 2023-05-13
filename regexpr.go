package govalidator

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type RegExprCheck struct {
	expr      *regexp.Regexp
	minLength *int
}

func (r *RegExprCheck) IsValidStringPtr(ctx context.Context, value interface{}) (bool, error) {
	if value == nil {
		return false, nil
	}
	actual, ok := value.(*string)
	if !ok {
		return false, fmt.Errorf("expected: %T, but had: %T", actual, value)
	}
	if actual == nil {
		return false, nil
	}
	if r.minLength != nil {
		if len(*actual) < *r.minLength {
			return false, nil
		}
	}
	return r.expr.MatchString(*actual), nil
}

func (r *RegExprCheck) IsValidString(ctx context.Context, value interface{}) (bool, error) {
	actual, ok := value.(string)
	if !ok {
		return false, fmt.Errorf("expected: %T, but had: %T", actual, value)
	}
	if r.minLength != nil {
		if len(actual) < *r.minLength {
			return false, nil
		}
	}
	return r.expr.MatchString(actual), nil
}

func NewNotRegExprCheck(expr *regexp.Regexp) func(field *Field, check *Check) (IsValid, error) {
	fn := NewRegExprCheck(expr)
	return func(field *Field, check *Check) (IsValid, error) {
		isValid, err := fn(field, check)
		if err != nil {
			return nil, err
		}
		return func(ctx context.Context, value interface{}) (bool, error) {
			ret, err := isValid(ctx, value)
			if err != nil {
				return false, err
			}
			return !ret, nil
		}, nil
	}
}

//NewRegExprCheck creates a regexpr based validation check
func NewRegExprCheck(expr *regexp.Regexp) func(field *Field, check *Check) (IsValid, error) {
	ret := &RegExprCheck{expr: expr}
	return func(field *Field, check *Check) (IsValid, error) {
		switch field.Kind() {
		case reflect.String:
			return ret.IsValidString, nil
		case reflect.Ptr:
			if field.Type.Elem().Kind() == reflect.String {
				return ret.IsValidStringPtr, nil
			}
		}
		return nil, fmt.Errorf("unsupported regexpr based %v check type: %s", field.Tag, field.Type.String())
	}
}

func NewRepeatedRegExprCheck(expr *regexp.Regexp, separator string) func(field *Field, check *Check) (IsValid, error) {
	fn := NewRegExprCheck(expr)
	return func(field *Field, check *Check) (IsValid, error) {
		isValid, err := fn(field, check)
		if err != nil {
			return nil, err
		}
		return func(ctx context.Context, value interface{}) (bool, error) {
			fragment := ""
			switch actual := value.(type) {
			case string:
				fragment = actual
			case []byte:
				fragment = string(actual)
			default:
				return false, fmt.Errorf("invalid input type: expected: %T, but had: %T", fragment, value)
			}
			if len(fragment) == 0 {
				return false, nil
			}
			for _, item := range strings.Split(fragment, separator) {
				ret, err := isValid(ctx, item)
				if err != nil {
					return false, err
				}
				if !ret {
					return false, nil
				}
			}
			return true, nil
		}, nil
	}
}

//NewRegExprCheckWithMinLength creates a regexpr based validation check with min Length
func NewRegExprCheckWithMinLength(expr *regexp.Regexp, minLen int) func(field *Field, check *Check) (IsValid, error) {
	ret := &RegExprCheck{expr: expr, minLength: &minLen}
	return func(field *Field, check *Check) (IsValid, error) {
		switch field.Kind() {
		case reflect.String:
			return ret.IsValidString, nil
		case reflect.Ptr:
			if field.Type.Elem().Kind() == reflect.String {
				return ret.IsValidStringPtr, nil
			}
		}
		return nil, fmt.Errorf("unsupported regexpr based %v check type: %s", field.Tag, field.Type.String())
	}
}
