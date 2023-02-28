package govalidator

import "context"

type (

	//Session represents validation session
	Session struct {
		Path        *Path
		Field       *Field
		ParentValue interface{}
	}
)

func (s *Session) Set(path *Path, field *Field, parentValue interface{}) {
	s.Path = path
	s.Field = field
	s.ParentValue = parentValue
}

//SessionKey represents a session key
var SessionKey string

//SessionContext creates a context with session
func SessionContext(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, SessionKey, session)
}
