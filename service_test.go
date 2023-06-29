package govalidator

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
		Has   *RecordHas `setMarker:"true"`
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
			description: "required string",
			input: struct {
				N *string `validate:"required"`
			}{N: stringPtr("ddd")},
			expectFailed: false,
		},
		{
			description: "required string",
			input: struct {
				N *string `validate:"required"`
			}{N: nil},
			expectFailed: true,
		},
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
					XX    int
					Email string `validate:"email"`
				}
			}{
				ID: 1,
				Contact: []struct {
					XX    int
					Email string `validate:"email"`
				}{{0, "xyz"}, {0, "zz@wp.pl"}, {0, "rrrr"}},
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
			description: "required string",
			input: struct {
				N string `validate:"required"`
			}{N: "ddd"},
			expectFailed: false,
		},
		{
			description: "required struct - zero value timestamp",
			input: struct {
				time time.Time `validate:"required"`
			}{time: time.Time{}}, // zero value time time.Parse(time.RFC3339, "0001-01-01T00:00:00Z")
			expectFailed: true,
		},
		{
			description: "required struct - non zero value timestamp",
			input: struct {
				time time.Time `validate:"required"`
			}{time: getTime("2023-06-29T23:09:15Z")},
			expectFailed: false,
		},
		{
			description: "With Presence pass",
			input: &Record{
				Id:  1,
				Has: &RecordHas{},
			},
			options:      []Option{WithSetMarker()},
			expectFailed: false,
		},
		{
			description: "With marker failed",
			input: &Record{
				Id: 1,
				Has: &RecordHas{
					Name: true,
				},
			},
			options:      []Option{WithSetMarker()},
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
				N     string
				F     float64
				Value int `validate:"gte=6"`
			}{"", 0.0, 5},
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
		{
			description: "choice valid",
			input: struct {
				Value string `validate:"choice(AZ,AK,ZZ),omitempty"`
			}{Value: "AK"},
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
			fmt.Printf("%v\n", validation.String())
			fmt.Printf("%v", validation)
		}
	}
}

func stringPtr(s string) *string {
	return &s
}

func getTime(timeStr string) time.Time {
	r, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}
	}
	return r
}
