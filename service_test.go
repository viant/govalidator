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

	type Req struct {
		Id int `validate:"required"`
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
			description: "required int",
			input: struct {
				N int `validate:"required"`
			}{N: 1},
			expectFailed: false,
		},
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
			description: "lte alias passed",
			input: struct {
				Value int `validate:"lte(5)"`
			}{Value: 5},
			expectFailed: false,
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
		{
			description: "required int64 supported",
			input: struct {
				Value int64 `validate:"required"`
			}{Value: 1},
			expectFailed: false,
		},
		{
			description: "required float32 supported",
			input: struct {
				Value float32 `validate:"required"`
			}{Value: 1.0},
			expectFailed: false,
		},
		{
			description: "rgba valid",
			input: struct {
				Value string `validate:"rgba"`
			}{Value: "rgba(255,255,255,0.1)"},
			expectFailed: false,
		},
		{
			description: "primitive int slice numeric checks item-wise",
			input: struct {
				Value []int `validate:"gt(3)"`
			}{Value: []int{1, 2, 4}},
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
			fmt.Printf("%v\n", validation.String())
			fmt.Printf("%v", validation)
		}
	}
}

func TestService_Validate_DoesNotPanicOnPointerString(t *testing.T) {
	value := "abc"
	srv := New()
	assert.NotPanics(t, func() {
		validation, err := srv.Validate(context.Background(), &value)
		assert.Nil(t, validation)
		assert.NotNil(t, err)
	})
}

func TestService_Validate_PrimitiveSliceNumericViolationCount(t *testing.T) {
	srv := New()
	validation, err := srv.Validate(context.Background(), struct {
		Value []int `validate:"gt(3)"`
	}{Value: []int{1, 2, 4}})
	if !assert.Nil(t, err) {
		return
	}
	if assert.True(t, validation.Failed) {
		assert.Equal(t, 2, len(validation.Violations))
	}
}

func TestService_Validate_NewCommonChecks(t *testing.T) {
	type Inner struct {
		ID int
	}
	tests := []struct {
		description  string
		input        interface{}
		expectFailed bool
		expectErr    bool
	}{
		{
			description: "min string length passes",
			input: struct {
				Name string `validate:"min(3)"`
			}{Name: "abcd"},
			expectFailed: false,
		},
		{
			description: "max string length fails",
			input: struct {
				Name string `validate:"max(3)"`
			}{Name: "abcd"},
			expectFailed: true,
		},
		{
			description: "between numeric passes",
			input: struct {
				Score float64 `validate:"between(1.5,2.5)"`
			}{Score: 2.0},
			expectFailed: false,
		},
		{
			description: "between numeric fails",
			input: struct {
				Score int `validate:"between(2,4)"`
			}{Score: 5},
			expectFailed: true,
		},
		{
			description: "min slice length passes",
			input: struct {
				Items []Inner `validate:"min(2)"`
			}{Items: []Inner{{ID: 1}, {ID: 2}}},
			expectFailed: false,
		},
		{
			description: "between primitive slice values passes item-wise",
			input: struct {
				Values []int `validate:"between(1,3)"`
			}{Values: []int{1, 2, 3}},
			expectFailed: false,
		},
		{
			description: "between primitive slice values fails item-wise",
			input: struct {
				Values []int `validate:"between(1,3)"`
			}{Values: []int{1, 4}},
			expectFailed: true,
		},
		{
			description: "oneof alias works",
			input: struct {
				State string `validate:"oneof(AZ,AK,ZZ)"`
			}{State: "AK"},
			expectFailed: false,
		},
		{
			description: "gte alias works",
			input: struct {
				Value int `validate:"gte(3)"`
			}{Value: 4},
			expectFailed: false,
		},
		{
			description: "lte alias works",
			input: struct {
				Value int `validate:"lte(3)"`
			}{Value: 2},
			expectFailed: false,
		},
		{
			description: "contains and startswith pass",
			input: struct {
				Value string `validate:"contains(foo),startswith(pre),endswith(post)"`
			}{Value: "pre-foo-post"},
			expectFailed: false,
		},
		{
			description: "notcontains fails",
			input: struct {
				Value string `validate:"notcontains(foo)"`
			}{Value: "foo-bar"},
			expectFailed: true,
		},
		{
			description: "eqfield passes",
			input: struct {
				Password        string `validate:"required"`
				ConfirmPassword string `validate:"eqfield(Password)"`
			}{Password: "x", ConfirmPassword: "x"},
			expectFailed: false,
		},
		{
			description: "nefield fails",
			input: struct {
				A int `validate:"required"`
				B int `validate:"nefield(A)"`
			}{A: 10, B: 10},
			expectFailed: true,
		},
		{
			description: "gtfield passes",
			input: struct {
				Min int `validate:"required"`
				Max int `validate:"gtfield(Min)"`
			}{Min: 5, Max: 7},
			expectFailed: false,
		},
		{
			description: "gtfield fails",
			input: struct {
				Min int `validate:"required"`
				Max int `validate:"gtfield(Min)"`
			}{Min: 5, Max: 4},
			expectFailed: true,
		},
		{
			description: "required_if enforced",
			input: struct {
				Type  string `validate:"required"`
				Phone string `validate:"required_if(Type,mobile)"`
			}{Type: "mobile", Phone: ""},
			expectFailed: true,
		},
		{
			description: "required_if skipped",
			input: struct {
				Type  string `validate:"required"`
				Phone string `validate:"required_if(Type,mobile)"`
			}{Type: "home", Phone: ""},
			expectFailed: false,
		},
		{
			description: "required_unless enforced",
			input: struct {
				Status string `validate:"required"`
				Reason string `validate:"required_unless(Status,ok)"`
			}{Status: "failed", Reason: ""},
			expectFailed: true,
		},
		{
			description: "required_with enforced",
			input: struct {
				Email string `validate:"omitempty,email"`
				Phone string `validate:"required_with(Email)"`
			}{Email: "test@example.com", Phone: ""},
			expectFailed: true,
		},
		{
			description: "required_without enforced",
			input: struct {
				Email string `validate:"omitempty,email"`
				Phone string `validate:"required_without(Email)"`
			}{Email: "", Phone: ""},
			expectFailed: true,
		},
		{
			description: "cross field missing target returns error",
			input: struct {
				First string `validate:"eqfield(Missing)"`
			}{First: "x"},
			expectErr: true,
		},
		{
			description: "url passes",
			input: struct {
				Value string `validate:"url"`
			}{Value: "https://example.com/path"},
			expectFailed: false,
		},
		{
			description: "url fails",
			input: struct {
				Value string `validate:"url"`
			}{Value: "example.com/path"},
			expectFailed: true,
		},
		{
			description: "uri passes",
			input: struct {
				Value string `validate:"uri"`
			}{Value: "urn:example:test"},
			expectFailed: false,
		},
		{
			description: "http_url passes",
			input: struct {
				Value string `validate:"http_url"`
			}{Value: "http://example.com"},
			expectFailed: false,
		},
		{
			description: "http_url fails for non-http scheme",
			input: struct {
				Value string `validate:"http_url"`
			}{Value: "ftp://example.com"},
			expectFailed: true,
		},
		{
			description: "ip passes",
			input: struct {
				Value string `validate:"ip"`
			}{Value: "2001:db8::1"},
			expectFailed: false,
		},
		{
			description: "ipv4 passes",
			input: struct {
				Value string `validate:"ipv4"`
			}{Value: "192.168.0.1"},
			expectFailed: false,
		},
		{
			description: "ipv6 passes",
			input: struct {
				Value string `validate:"ipv6"`
			}{Value: "2001:db8::1"},
			expectFailed: false,
		},
		{
			description: "cidr passes",
			input: struct {
				Value string `validate:"cidr"`
			}{Value: "10.0.0.0/24"},
			expectFailed: false,
		},
		{
			description: "hostname passes",
			input: struct {
				Value string `validate:"hostname"`
			}{Value: "api.example-1.com"},
			expectFailed: false,
		},
		{
			description: "mac passes",
			input: struct {
				Value string `validate:"mac"`
			}{Value: "01:23:45:67:89:ab"},
			expectFailed: false,
		},
		{
			description: "port passes for string",
			input: struct {
				Value string `validate:"port"`
			}{Value: "443"},
			expectFailed: false,
		},
		{
			description: "port passes for int",
			input: struct {
				Value int `validate:"port"`
			}{Value: 8080},
			expectFailed: false,
		},
		{
			description: "port fails for out of range",
			input: struct {
				Value int `validate:"port"`
			}{Value: 70000},
			expectFailed: true,
		},
		{
			description: "uuidv7 passes",
			input: struct {
				Value string `validate:"uuidv7"`
			}{Value: "01890f47-8c98-7f4a-8af5-67c4c8602f7a"},
			expectFailed: false,
		},
		{
			description: "uuidv7 fails for version mismatch",
			input: struct {
				Value string `validate:"uuidv7"`
			}{Value: "01890f47-8c98-4f4a-8af5-67c4c8602f7a"},
			expectFailed: true,
		},
		{
			description: "slug passes",
			input: struct {
				Value string `validate:"slug"`
			}{Value: "hello-world-123"},
			expectFailed: false,
		},
		{
			description: "slug fails",
			input: struct {
				Value string `validate:"slug"`
			}{Value: "Hello World"},
			expectFailed: true,
		},
		{
			description: "semver passes",
			input: struct {
				Value string `validate:"semver"`
			}{Value: "1.2.3-beta.1+build.10"},
			expectFailed: false,
		},
		{
			description: "semver fails",
			input: struct {
				Value string `validate:"semver"`
			}{Value: "1.2"},
			expectFailed: true,
		},
		{
			description: "json passes for string",
			input: struct {
				Value string `validate:"json"`
			}{Value: "{\"a\":1}"},
			expectFailed: false,
		},
		{
			description: "json fails for string",
			input: struct {
				Value string `validate:"json"`
			}{Value: "{a:1}"},
			expectFailed: true,
		},
		{
			description: "past passes for time",
			input: struct {
				Value time.Time `validate:"past"`
			}{Value: time.Now().Add(-1 * time.Hour)},
			expectFailed: false,
		},
		{
			description: "future fails for past time",
			input: struct {
				Value time.Time `validate:"future"`
			}{Value: time.Now().Add(-1 * time.Hour)},
			expectFailed: true,
		},
		{
			description: "future passes for RFC3339 string",
			input: struct {
				Value string `validate:"future"`
			}{Value: time.Now().Add(2 * time.Hour).Format(time.RFC3339)},
			expectFailed: false,
		},
		{
			description: "past fails for RFC3339 string",
			input: struct {
				Value string `validate:"past"`
			}{Value: time.Now().Add(2 * time.Hour).Format(time.RFC3339)},
			expectFailed: true,
		},
	}

	for _, test := range tests {
		validation, err := New().Validate(context.Background(), test.input)
		if test.expectErr {
			assert.NotNil(t, err, test.description)
			continue
		}
		if !assert.Nil(t, err, test.description) {
			continue
		}
		assert.EqualValues(t, test.expectFailed, validation.Failed, test.description)
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
