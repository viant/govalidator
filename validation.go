package govalidator

import (
	"fmt"
	"github.com/viant/xunsafe"
	"reflect"
	"strings"
	"unsafe"
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

func (v *Validation) AddViolation(field string, value interface{}, check string, msg string) {
	path := &Path{Kind: PathKinField, Name: field}
	v.Append(path, field, value, check, msg)
}

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
	v.Failed = len(v.Violations) > 0
}

func (e *Validation) Error() string {
	return e.String()
}

func (e *Validation) String() string {
	if e == nil || len(e.Violations) == 0 {
		return ""
	}
	msg := strings.Builder{}
	msg.WriteString("Failed validation for ")
	for i, v := range e.Violations {
		if i > 0 {
			msg.WriteString(",")
		}
		msg.WriteString(v.Location)
		msg.WriteString("(")
		msg.WriteString(v.Check)
		msg.WriteString(")")

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
	v := reflect.TypeOf(value)
	if v.Kind() == reflect.Ptr {
		ptr := xunsafe.AsPointer(value)
		if ptr == nil || (*unsafe.Pointer)(ptr) == nil {
			return value
		}
	}
	return value
}
