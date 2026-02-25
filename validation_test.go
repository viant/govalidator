package govalidator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidation_AppendTemplateVariables(t *testing.T) {
	validation := &Validation{}
	path := NewPath().Field("Confirm")
	validation.Append(path, "Confirm", "y", "eqfield", "$field:$value:$param:$otherField", []string{"Password"})
	if assert.Equal(t, 1, len(validation.Violations)) {
		assert.Equal(t, "Confirm:y:Password:Password", validation.Violations[0].Message)
	}
}
