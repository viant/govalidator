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

		{
			description: "valid domain",
			input: struct {
				Phone *string `validate:"omitempty,domain"`
			}{Phone: stringPtr("wp.pl")},
			expectFailed: false,
		},
		{
			description: "invalid domain",
			input: struct {
				Phone *string `validate:"omitempty,domain"`
			}{Phone: stringPtr("wp-.pl")},
			expectFailed: true,
		},
		{
			description: "valid www domain",
			input: struct {
				Value *string `validate:"omitempty,wwwDomain"`
			}{Value: stringPtr("www.wp.pl")},
			expectFailed: false,
		},
		{
			description: "invalid www domain",
			input: struct {
				Value *string `validate:"omitempty,wwwDomain"`
			}{Value: stringPtr("lll.wp.pl")},
			expectFailed: true,
		},

		{
			description: "valid top level domain",
			input: struct {
				Value *string `validate:"omitempty,domain,nonWWWDomain"`
			}{Value: stringPtr("lll.wp.pl")},
			expectFailed: false,
		},
		{
			description: "invalid top level domain",
			input: struct {
				Value *string `validate:"omitempty,domain,nonWWWDomain"`
			}{Value: stringPtr("www.wp.pl")},
			expectFailed: true,
		},

		{
			description: "less than 3 - violation",
			input: struct {
				Value *string `validate:"omitempty,lt(3)"`
			}{Value: stringPtr("434")},
			expectFailed: true,
		},
		{
			description: "less than 3 - valid",
			input: struct {
				Value string `validate:"omitempty,lt(5)"`
			}{Value: "123"},
			expectFailed: false,
		},
		{
			description: "less than 300- valid",
			input: struct {
				Value int `validate:"omitempty,gt(0),lt(300)"`
			}{Value: 200},
			expectFailed: false,
		},
		{
			description: "less between 0 .. 300 - invalid",
			input: struct {
				Value int `validate:"omitempty,gt(0),lt(100)"`
			}{Value: 200},
			expectFailed: true,
		},
		{
			description: "less between 0 .. 300 - invalid",
			input: struct {
				Value int `validate:"omitempty,gt(0),lt(100)"`
			}{Value: -1},
			expectFailed: true,
		},
		{
			description: "IAB category - valid 1",
			input: struct {
				Value string `validate:"omitempty,iabCategory"`
			}{Value: "IAB2-22"},
			expectFailed: false,
		},
		{
			description: "IAB category- valid 2",
			input: struct {
				Value string `validate:"omitempty,iabCategory"`
			}{Value: "IAB2"},
			expectFailed: false,
		},
		{
			description: "IAB categories - valid",
			input: struct {
				Value string `validate:"omitempty,iabCategories"`
			}{Value: "IAB2-22,IAB8"},
			expectFailed: false,
		},
		{
			description: "IAB categories - invalid",
			input: struct {
				Value string `validate:"omitempty,iabCategories"`
			}{Value: "sAB2-22,IAB8"},
			expectFailed: true,
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
