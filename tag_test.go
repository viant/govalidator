package govalidator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTag(t *testing.T) {

	var testCases = []struct {
		description string
		tag         string
		expect      *Tag
	}{
		{
			description: "required tag",
			tag:         "required",
			expect:      &Tag{Required: true, Checks: []Check{{Name: "required", Parameters: emptyArgs}}},
		},
		{
			description: "multi tags",
			tag:         "required|check1",
			expect:      &Tag{Required: true, Checks: []Check{{Name: "required", Parameters: emptyArgs}, {Name: "check1", Parameters: emptyArgs}}},
		},
		{
			description: "allow empy check with parameter",
			tag:         "omitempty|checkX(param1, param2)",
			expect:      &Tag{Omitempty: true, Checks: []Check{{Name: "checkX", Parameters: []string{"param1", "param2"}}}},
		},
	}

	for _, testCase := range testCases {
		actual := ParseTag(testCase.tag)
		assert.EqualValues(t, testCase.expect, actual, testCase.description)
	}

}
