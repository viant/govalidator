package govalidator

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
)

type Choice struct {
	intValue     map[int]bool
	stringValues map[string]bool
}

func (c *Choice) checkInts(ctx context.Context, value interface{}) (bool, error) {
	if isBaseTypeNil(value) {
		return false, nil
	}
	key := intValue(value)
	return c.intValue[key], nil
}

func (c *Choice) checkStrings(ctx context.Context, value interface{}) (bool, error) {
	if isBaseTypeNil(value) {
		return false, nil
	}
	key := stringValue(value)
	return c.stringValues[key], nil
}

func (c *Choice) setInts(args []string) error {
	c.intValue = make(map[int]bool)
	for _, arg := range args {
		key, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("%w invalid choice option expect int, but had: %v", err, arg)
		}
		c.intValue[key] = true
	}
	return nil
}

func (c *Choice) setStrings(args []string) {
	c.stringValues = make(map[string]bool)
	for i := range args {
		c.stringValues[args[i]] = true
	}

}

//NewChoice creates choice/enum value checks
func NewChoice() func(field *Field, check *Check) (IsValid, error) {

	return func(field *Field, check *Check) (IsValid, error) {
		choice := &Choice{}

		if len(check.Parameters) == 0 {
			return nil, fmt.Errorf("check option was empty")
		}

		switch field.Kind() {
		case reflect.Int, reflect.Uint, reflect.Int64, reflect.Uint64:
			if err := choice.setInts(check.Parameters); err != nil {
				return nil, err
			}
			return choice.checkInts, nil
		case reflect.String:
			choice.setStrings(check.Parameters)
			return choice.checkStrings, nil
		case reflect.Ptr:
			switch field.Elem().Kind() {
			case reflect.Int, reflect.Uint, reflect.Int64, reflect.Uint64:
				if err := choice.setInts(check.Parameters); err != nil {
					return nil, err
				}
				return choice.checkInts, nil
			case reflect.String:
				choice.setStrings(check.Parameters)
				return choice.checkStrings, nil
			}
		}
		return nil, fmt.Errorf("not supported for: %v", field.Type.String())
	}
}
