package govalidator

import (
	"context"
	"fmt"
	"github.com/viant/xunsafe"
	"reflect"
	"sync"
	"unsafe"
)

// Service represents a service
type Service struct {
	checks map[reflect.Type]*Checks
	mux    sync.RWMutex
}

func (s *Service) Validate(ctx context.Context, any interface{}, opts ...Option) (*Validation, error) {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}
	rootPath := NewPath()
	if options.Path != nil {
		rootPath.Path = options.Path
	}
	validation := &Validation{}
	ctx = SessionContext(ctx, &Session{Path: rootPath})

	if err := s.validate(ctx, any, validation, options); err != nil {
		return nil, err
	}
	validation.Failed = len(validation.Violations) > 0
	return validation, nil
}

func (s *Service) validate(ctx context.Context, any interface{}, validation *Validation, options *Options) error {
	if any == nil {
		return nil
	}
	t := reflect.TypeOf(any)
	switch t.Kind() {
	case reflect.Ptr:
		any = deref(any)
		switch t.Elem().Kind() {
		case reflect.Struct:
			return s.validateStruct(ctx, t.Elem(), any, validation, options)
		case reflect.Slice:
			return s.validateSlice(ctx, t.Elem(), any, validation, options)
		case reflect.String:
			actual := any.([]string)
			for _, value := range actual {
				if err := s.validate(ctx, value, validation, options); err != nil {
					return err
				}
			}
		}
	case reflect.Struct:
		return s.validateStruct(ctx, t, any, validation, options)
	case reflect.Slice:
		return s.validateSlice(ctx, t, any, validation, options)
	}
	return fmt.Errorf("unsupported value: %v", any)
}

func deref(any interface{}) interface{} {
	if any == nil {
		return nil
	}
	if zeroer, ok := any.(Zeroable); ok {
		if zeroer.IsZero() {
			return nil
		}
	}
	ptr := xunsafe.AsPointer(any)
	if ptr == nil {
		return nil
	}
	any = reflect.ValueOf(any).Elem().Interface()
	return any
}

func (s *Service) validateStruct(ctx context.Context, t reflect.Type, value interface{}, validation *Validation, options *Options) error {
	if value == nil {
		return nil
	}
	session := ctx.Value(SessionKey).(*Session)
	ptr := xunsafe.AsPointer(value)
	checks, err := s.checksFor(t)
	if err != nil {
		return err
	}
	path, field, parentValue := session.Path, session.Field, session.ParentValue
	defer session.Set(path, field, parentValue)
	if err := s.checkStructFields(ctx, checks, path, ptr, session, value, validation, options); err != nil {
		return err
	}
	if options.Shallow {
		return nil
	}
	err = s.diveStructFields(ctx, checks, path, ptr, session, validation, options)
	if err != nil {
		return err
	}
	err = s.diveSliceFields(ctx, checks, path, ptr, session, validation, options)
	if err != nil {
		return err
	}
	err = s.diveSimpleSliceFields(ctx, checks, path, ptr, session, validation, options)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) diveStructFields(ctx context.Context, checks *Checks, path *Path, ptr unsafe.Pointer, session *Session, validation *Validation, options *Options) error {
	if len(checks.Structs) == 0 {
		return nil
	}
	for _, candidate := range checks.Structs {
		fieldPath := path.Field(candidate.Name)
		if candidate.SkipPath {
			fieldPath = path
		}
		fieldValue := candidate.Value(ptr)
		if fieldValue == nil {
			continue
		}
		if candidate.Kind() == reflect.Ptr {
			fieldValue = deref(fieldValue)
		}
		session.Set(fieldPath, candidate, fieldValue)
		if err := s.validateStruct(ctx, candidate.Type, fieldValue, validation, options); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) diveSliceFields(ctx context.Context, checks *Checks, path *Path, ptr unsafe.Pointer, session *Session, validation *Validation, options *Options) error {

	if len(checks.Slices) == 0 {
		return nil
	}
	for _, candidate := range checks.Slices {
		fieldPath := path.Field(candidate.Name)
		if candidate.SkipPath {
			fieldPath = path
		}
		fieldValue := candidate.Value(ptr)
		if fieldValue == nil {
			continue
		}
		if candidate.Kind() == reflect.Ptr {
			fieldValue = deref(fieldValue)
		}
		session.Set(fieldPath, candidate, fieldValue)
		if err := s.validateSlice(ctx, candidate.Type, fieldValue, validation, options); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) diveSimpleSliceFields(ctx context.Context, checks *Checks, path *Path, ptr unsafe.Pointer, session *Session, validation *Validation, options *Options) error {

	if len(checks.SimpleSlices) == 0 {
		return nil
	}
	for _, candidate := range checks.SimpleSlices {
		fieldPath := path.Field(candidate.Name)
		if candidate.SkipPath {
			fieldPath = path
		}

		fieldValue := candidate.Value(ptr)
		if fieldValue == nil {
			continue
		}
		if candidate.Kind() == reflect.Ptr {
			fieldValue = deref(fieldValue)
		}
		if candidate.FieldCheck == nil {
			continue
		}
		fieldCheck := candidate.FieldCheck
		session.Set(fieldPath, candidate, fieldValue)
		switch actual := fieldValue.(type) {
		case []string:
			for j, item := range actual {
				elemPath := fieldPath.Element(j)
				session.Set(elemPath, candidate, item)
				if err := s.checkValue(ctx, fieldCheck, item, options, validation, elemPath); err != nil {
					return err
				}
			}
		case []int:
			for j, item := range actual {
				elemPath := fieldPath.Element(j)
				session.Set(elemPath, candidate, item)
				if err := s.checkValue(ctx, fieldCheck, item, options, validation, elemPath); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *Service) checkStructFields(ctx context.Context, checks *Checks, path *Path, ptr unsafe.Pointer, session *Session, value interface{}, validation *Validation, options *Options) error {
	if len(checks.Fields) == 0 {
		return nil
	}
	presence := checks.marker

	canUseMarker := options.UseMarker && presence.CanUseHolder(ptr)
	if options.CanUseMarkerProvider != nil && canUseMarker {
		canUseMarker = options.CanUseMarkerProvider(value)
	}

	for _, field := range checks.Fields {
		fieldPath := path.Field(field.Field.Name)
		fieldValue := field.Field.Value(ptr)

		if isEmpty(fieldValue) && field.Omitempty {
			continue
		}

		if options.UseMarker && presence.CanUseHolder(ptr) {
			if !presence.IsSet(ptr, int(field.Field.Index)) {
				continue
			}
		}

		session.Set(path, field.Field, value)
		err := s.checkValue(ctx, field, fieldValue, options, validation, fieldPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) checkValue(ctx context.Context, field *FieldCheck, fieldValue interface{}, options *Options, validation *Validation, fieldPath *Path) error {
	for i, isValid := range field.IsValid {
		passed, err := isValid(ctx, fieldValue)
		if err != nil {
			return err
		}
		if !passed {
			if field.Type.Kind() == reflect.Ptr && !options.PreservePointer {
				if isNil := isNil(fieldValue); isNil {
					fieldValue = nil
				} else {
					fieldValue = deref(fieldValue)
				}
			}
			validation.Append(fieldPath, field.Field.Name, fieldValue, field.Checks[i].Name, field.Checks[i].Message)
			break
		}
	}
	return nil
}

func (s *Service) validateSlice(ctx context.Context, t reflect.Type, any interface{}, validation *Validation, options *Options) error {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	xSlice := xunsafe.NewSlice(t)
	session := ctx.Value(SessionKey).(*Session)
	slicePtr := xunsafe.AsPointer(any)
	sliceLen := xSlice.Len(slicePtr)

	path, field, parentValue := session.Path, session.Field, session.ParentValue
	defer session.Set(path, field, parentValue)

	for i := 0; i < sliceLen; i++ {
		itmPath := path.Element(i)
		value := xSlice.ValueAt(slicePtr, i)
		session.Set(itmPath, field, any)
		if err := s.validateStruct(ctx, t.Elem(), value, validation, options); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) checksFor(t reflect.Type) (*Checks, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	s.mux.RLock()
	checks, ok := s.checks[t]
	s.mux.RUnlock()
	if ok {
		return checks, nil
	}
	var err error
	if checks, err = NewChecks(t); err != nil {
		return nil, err
	}
	s.mux.Lock()
	s.checks[t] = checks
	s.mux.Unlock()
	return checks, nil
}

func isEmpty(value interface{}) bool {
	switch actual := value.(type) {
	case *int, *uint, *int64, *uint64:
		ptr := (*int)(xunsafe.AsPointer(actual))
		return ptr == nil
	case *uint8:
		return actual == nil
	case *string:
		if actual == nil {
			return true
		}
		return *actual == ""
	case string:
		return actual == ""
	default:
		if value == nil {
			return true
		}
		if zeroer, ok := value.(Zeroable); ok {
			return zeroer.IsZero()
		}
		return value == nil
	}
}

func isNil(value interface{}) bool {
	switch actual := value.(type) {
	case *int, *uint, *int64, *uint64:
		ptr := (*int)(xunsafe.AsPointer(actual))
		return ptr == nil
	case *uint8:
		return actual == nil
	case *string:
		if actual == nil {
			return true
		}
	default:
		if value == nil {
			return true
		}
		if zeroer, ok := value.(Zeroable); ok {
			return zeroer.IsZero()
		}
		return value == nil
	}
	return false
}

func New() *Service {
	return &Service{
		checks: map[reflect.Type]*Checks{},
		mux:    sync.RWMutex{},
	}
}
