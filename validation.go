package govalidator

import (
	"fmt"
	"reflect"
	"strings"
)

type (
	Violation struct {
		Location string
		Field    string
		Value    interface{}
		Message  string
		Check    string
	}

	Validation struct {
		Violations []*Violation
		Failed     bool
	}
)

func (v *Validation) Append(path *Path, field string, value interface{}, check string, msg string) {
	value = derefIfNeeded(value)
	if msg == "" {
		msg = fmt.Sprintf("check '%v' failed on field %v", check, field)
	} else {
		msg = strings.Replace(msg, "$value", fmt.Sprintf("%v", value), 1)
	}
	v.Violations = append(v.Violations, &Violation{
		Location: path.String(),
		Field:    field,
		Message:  msg,
		Check:    check,
		Value:    value,
	})
}

func (e *Validation) Error() string {
	return e.String()
}

func (e *Validation) String() string {
	if e == nil || len(e.Violations) == 0 {
		return ""
	}
	msg := strings.Builder{}
	msg.WriteString("Field validation for ")
	for i, v := range e.Violations {
		if i > 0 {
			msg.WriteString(",")
		}
		msg.WriteString(v.Location)
	}
	return msg.String()
}

func derefIfNeeded(value interface{}) interface{} {
	switch actual := value.(type) {
	case *string:
		if actual == nil {
			return nil
		}
	case *int, *int64, *uint, *uint64:
		if actual == nil {
			return nil
		}
	case *float32, *float64:
		if actual == nil {
			return nil
		}
	}
	if value == nil {
		return nil
	}
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		value = v.Interface()
	}
	return value
}
