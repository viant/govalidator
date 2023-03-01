package govalidator

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_Validate(t *testing.T) {

	type RecordHas struct {
		Id    bool
		Name  bool
		Phone bool
		Email bool
	}

	type Record struct {
		Id    int
		Name  *string    `validate:"required"`
		Phone string     `validate:"phone"`
		Email string     `validate:"omitempty,email"`
		Has   *RecordHas `validate:"presence"`
	}

	type BasicRecord struct {
		ID    int
		Phone string `validate:"phone"`
		Email string `validate:"omitempty,email"`
	}

	type RequiredCheck struct {
		ID   int
		Name string `validate:"required"`
	}

	type SkipRoot struct {
		ID     int
		Record *BasicRecord `validate:"skipPath"`
	}

	var testCases = []struct {
		description  string
		input        interface{}
		expectFailed bool
		options      []Option
	}{
		{
			description: "basic email validation",
			input: struct {
				Email string `validate:"email"`
			}{
				Email: "abc",
			},
			expectFailed: true,
		},
		{
			description: "nested validation",
			input: struct {
				ID      int
				Contact struct {
					Email string `validate:"email"`
				}
			}{
				ID: 1,
				Contact: struct {
					Email string `validate:"email"`
				}{"xyz"},
			},
			expectFailed: true,
		},
		{
			description: "repeated validation",
			input: struct {
				ID      int
				Contact []struct {
					Email string `validate:"email"`
				}
			}{
				ID: 1,
				Contact: []struct {
					Email string `validate:"email"`
				}{{"xyz"}, {"zz@wp.pl"}, {"rrrr"}},
			},
			expectFailed: true,
		},
		{
			description:  "repeated ptr validation",
			input:        []*BasicRecord{{ID: 1, Phone: "213-300-2222"}, {ID: 1, Phone: "213-300-22222"}, {ID: 1, Email: "aaa"}},
			expectFailed: true,
		},
		{
			description:  "valid phone",
			input:        &BasicRecord{ID: 1, Phone: "213-300-2222"},
			expectFailed: false,
		},
		{
			description:  "repeated ptr validation",
			input:        SkipRoot{Record: &BasicRecord{ID: 1, Phone: "213-300-085"}},
			expectFailed: true,
		},
		{
			description:  "required",
			input:        &RequiredCheck{ID: 1323},
			expectFailed: true,
		},
		{
			description: "With Presence pass",
			input: &Record{
				Id:  1,
				Has: &RecordHas{},
			},
			options:      []Option{WithPresence()},
			expectFailed: false,
		},
		{
			description: "With presence failed",
			input: &Record{
				Id: 1,
				Has: &RecordHas{
					Name: true,
				},
			},
			options:      []Option{WithPresence()},
			expectFailed: true,
		},
		{
			description:  "shallow",
			input:        SkipRoot{Record: &BasicRecord{ID: 1, Phone: "213-300-085"}},
			expectFailed: false,
			options:      []Option{WithShallow(true)},
		},
		{
			description: "ge passed",
			input: struct {
				Value int `validate:"ge=3"`
			}{5},
			expectFailed: false,
		},
		{
			description: "gte failed",
			input: struct {
				Value int `validate:"gte=6"`
			}{5},
			expectFailed: true,
		},
		{
			description: "phone valid ptr",
			input: struct {
				Phone *string `validate:"omitempty,phone"`
			}{Phone: stringPtr("213-222-0001")},
			expectFailed: false,
		},
	}

	for _, testCase := range testCases {
		srv := New()
		validation, err := srv.Validate(context.Background(), testCase.input, testCase.options...)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}
		if !assert.EqualValues(t, testCase.expectFailed, validation.Failed, testCase.description) {
			fmt.Printf("%v", validation)
		}
	}
}

func stringPtr(s string) *string {
	return &s
}
