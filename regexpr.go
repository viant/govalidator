package govalidator

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
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
