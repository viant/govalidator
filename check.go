package govalidator

import (
	"fmt"
	"github.com/viant/structology"
	"github.com/viant/xunsafe"
	"reflect"
	"time"
)

type (

	//FieldCheckPos represents field checks
	FieldCheck struct {
		*Field
		Owner   reflect.Type
		IsValid []IsValid
	}

	Field struct {
		*Tag
		*xunsafe.Field
		FieldCheck *FieldCheck
	}

	//Checks represents struct checks
	Checks struct {
		Type         reflect.Type
		Fields       []*FieldCheck
		Slices       []*Field
		Structs      []*Field
		SimpleSlices []*Field
		marker       *structology.Marker
	}
)

// NewChecks returns new checks
func NewChecks(t reflect.Type) (*Checks, error) {
	checks := &Checks{Type: t}
	sType := t
	if sType.Kind() == reflect.Ptr {
		sType = sType.Elem()
	}
	checks.marker, _ = structology.NewMarker(sType)
	xStruct := xunsafe.NewStruct(sType)

	var fieldPos = map[string]int{}
	for i := range xStruct.Fields {
		xField := &xStruct.Fields[i]
		fieldPos[xField.Name] = int(xField.Index)
		tagLiteral, ok := xField.Tag.Lookup("validate")
		tag := ParseTag(tagLiteral)
		if structology.IsSetMarker(xField.Tag) {
			continue
		}
		if isStruct(xField.Type) && !isTime(xField.Type) {
			checks.Structs = append(checks.Structs, &Field{Tag: tag, Field: xField})
		} else if isSliceStruct(xField.Type) {
			checks.Slices = append(checks.Slices, &Field{Tag: tag, Field: xField})
		} else if xField.Type.Kind() == reflect.Slice && isPrimitive(xField.Type.Elem()) {
			field := &Field{Tag: tag, Field: xField}
			if ok {
				fieldCheck, err := buildFieldCheck(sType, field, tag)
				if err != nil {
					return nil, err
				}
				field.FieldCheck = fieldCheck
			}
			checks.SimpleSlices = append(checks.SimpleSlices, field)
			continue
		}
		if !ok || tagLiteral == "" {
			continue
		}
		field := &Field{Tag: tag, Field: xField}
		fieldCheck, err := buildFieldCheck(sType, field, tag)
		if err != nil {
			return nil, err
		}
		checks.Fields = append(checks.Fields, fieldCheck)
	}
	return checks, nil
}

func buildFieldCheck(sType reflect.Type, field *Field, tag *Tag) (*FieldCheck, error) {
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
	return fieldCheck, nil
}

func isPrimitive(t reflect.Type) bool {
	if t.Kind() == reflect.String {
		return true
	}
	if t.Kind() == reflect.Int {
		return true
	}
	return false
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
