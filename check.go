package govalidator

import (
	"fmt"
	"github.com/viant/xunsafe"
	"reflect"
	"time"
)

type (

	//FieldCheck represents field checks
	FieldCheck struct {
		*Field
		Owner   reflect.Type
		IsValid []IsValid
	}

	Field struct {
		*Tag
		*xunsafe.Field
	}

	//Checks represents struct checks
	Checks struct {
		Type   reflect.Type
		Fields []*FieldCheck

		Slices  []*Field
		Structs []*Field

		presence *PresenceProvider
	}
)

//NewChecks returns new checks
func NewChecks(t reflect.Type, presence *PresenceProvider) (*Checks, error) {
	checks := &Checks{Type: t, presence: presence}
	sType := t
	if sType.Kind() == reflect.Ptr {
		sType = sType.Elem()
	}
	xStruct := xunsafe.NewStruct(sType)
	var fieldPos = map[string]int{}
	for i := range xStruct.Fields {
		xField := &xStruct.Fields[i]
		fieldPos[xField.Name] = int(xField.Index)
		tagLiteral, ok := xField.Tag.Lookup("validate")
		tag := ParseTag(tagLiteral)
		if !tag.Presence {
			if _, ok := xField.Tag.Lookup("presenceIndex"); ok {
				tag.Presence = ok
			}
		}
		if presence != nil && tag.Presence {
			presence.Holder = xField
		}

		if isStruct(xField.Type) && !isTime(xField.Type) {
			checks.Structs = append(checks.Structs, &Field{Tag: tag, Field: xField})
		} else if isSliceStruct(xField.Type) {
			checks.Slices = append(checks.Slices, &Field{Tag: tag, Field: xField})
		}
		if !ok || tagLiteral == "" {
			continue
		}
		field := &Field{Tag: tag, Field: xField}
		fieldCheck := &FieldCheck{Owner: sType, Field: field}
		for i := range tag.Checks {
			check := &tag.Checks[i]
			newCheck := LookupAll(check.Name)
			if newCheck == nil {
				return nil, fmt.Errorf("unknown check: %v", check.Name)
			}
			isValid, err := newCheck(field, check)
			if err != nil {
				return nil, err
			}
			fieldCheck.IsValid = append(fieldCheck.IsValid, isValid)

		}
		checks.Fields = append(checks.Fields, fieldCheck)
	}

	if presence != nil && presence.Holder != nil {
		if err := presence.Init(fieldPos); err != nil {
			return nil, err
		}
	}
	return checks, nil
}

func isSliceStruct(t reflect.Type) bool {
	if t.Kind() == reflect.Slice {
		return isStruct(t.Elem())
	}
	if t.Kind() == reflect.Ptr {
		return isSliceStruct(t.Elem())
	}
	return false
}

func isStruct(t reflect.Type) bool {
	if t.Kind() == reflect.Struct {
		return true
	}
	if t.Kind() == reflect.Ptr {
		return isStruct(t.Elem())
	}
	return false
}

func isTime(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t == reflect.TypeOf(time.Time{})
}
