package govalidator

import "context"

//IsValid represents a validation function
type IsValid func(ctx context.Context, value interface{}) (bool, error)

//NewIsValid function to create IsValid
type NewIsValid func(field *Field, check *Check) (IsValid, error)

//atListOneValid return aliased valid check
func atListOneValid(group ...NewIsValid) func(field *Field, check *Check) (IsValid, error) {
	return func(field *Field, check *Check) (IsValid, error) {
		var validFns = make([]IsValid, 0)
		for i := range group {
			isValid, err := group[i](field, check)
			if err != nil {
				return nil, err
			}
			validFns = append(validFns, isValid)
		}
		return func(ctx context.Context, value interface{}) (bool, error) {
			for _, isValid := range validFns {
				passed, err := isValid(ctx, value)
				if err != nil || passed {
					return passed, err
				}
			}
			return false, nil
		}, nil
	}
}
