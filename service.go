package govalidator

import (
	"context"
	"fmt"
	"github.com/viant/xunsafe"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

//Service represents a service
type Service struct {
	checks map[reflect.Type]*Checks
	mux    sync.RWMutex
}

func (s *Service) Validate(ctx context.Context, any interface{}, opts ...Option) (*Validation, error) {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}
	validation := &Validation{}
	ctx = SessionContext(ctx, &Session{Path: &Path{}})
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
		}
	case reflect.Struct:
		return s.validateStruct(ctx, t, any, validation, options)
	case reflect.Slice:
		return s.validateSlice(ctx, t, any, validation, options)
	}
	return fmt.Errorf("unsupported value: %v", any)
}

func deref(any interface{}) interface{} {
	any = reflect.ValueOf(any).Elem().Interface()
	return any
}

func (s *Service) validateStruct(ctx context.Context, t reflect.Type, value interface{}, validation *Validation, options *Options) error {
	session := ctx.Value(SessionKey).(*Session)
	ptr := xunsafe.AsPointer(value)
	checks, err := s.checksFor(t, options.PresenceProvider)
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

func (s *Service) checkStructFields(ctx context.Context, checks *Checks, path *Path, ptr unsafe.Pointer, session *Session, value interface{}, validation *Validation, options *Options) error {
	if len(checks.Fields) == 0 {
		return nil
	}
	presence := options.PresenceProvider
	for i, field := range checks.Fields {
		fieldPath := path.Field(field.Field.Name)
		fieldValue := field.Field.Value(ptr)

		if isEmpty(fieldValue) && field.Omitempty {
			continue
		}
		if !presence.IsFieldSet(ptr, i) {
			continue
		}

		session.Set(path, field.Field, value)
		for i, isValid := range field.IsValid {
			passed, err := isValid(ctx, fieldValue)
			if err != nil {
				return err
			}
			if !passed {
				validation.Append(fieldPath, field.Field.Name, fieldValue, field.Checks[i].Name, field.Checks[i].Message)
				break
			}
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

func (s *Service) checksFor(t reflect.Type, presence *PresenceProvider) (*Checks, error) {
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
	if checks, err = NewChecks(t, presence); err != nil {
		return nil, err
	}
	s.mux.Lock()
	s.checks[t] = checks
	s.mux.Unlock()
	return checks, nil
}

func isNil(value interface{}) bool {
	switch actual := value.(type) {
	case *int, *uint, *int64, *uint64:
		ptr := (*int)(xunsafe.AsPointer(actual))
		return ptr == nil
	case *uint8:
		return actual == nil
	case *string:
		return actual == nil
	case *time.Time:
		return actual == nil
	default:
		return value == nil
	}
}
func isEmpty(value interface{}) bool {
	switch actual := value.(type) {
	case *int, *uint, *int64, *uint64:
		ptr := (*int)(xunsafe.AsPointer(actual))
		return ptr == nil
	case *uint8:
		return actual == nil
	case *string:
		return actual == nil || *actual == ""
	case string:
		return actual == ""
	default:
		if value == nil {
			return true
		}
		if zeroer, ok := value.(IsZero); ok {
			return zeroer.IsZero()
		}
		return value == nil
	}
}

func New() *Service {
	return &Service{
		checks: map[reflect.Type]*Checks{},
		mux:    sync.RWMutex{},
	}
}
