package validator

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator_Check(t *testing.T) {
	tests := []struct {
		name        string
		ok          bool
		key         string
		message     string
		expectError bool
	}{
		{
			name:        "Should not add error when check passes",
			ok:          true,
			key:         "ErrorKey",
			message:     "ErrorMessage",
			expectError: false,
		},
		{
			name:        "Should add error when check doesn't pass",
			ok:          false,
			key:         "ErrorKey",
			message:     "ErrorMessage",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			validator := New()
			validator.Check(tc.ok, tc.key, tc.message)
			errors := validator.Errors()

			if tc.expectError {
				assert.Equal(t, 1, len(errors))
				value, ok := errors[tc.key]
				assert.True(t, ok)
				assert.Equal(t, tc.message, value)
			} else {
				assert.Equal(t, 0, len(errors))
			}
		})
	}
}

func TestValidator_Check_SameError(t *testing.T) {
	validator := New()
	validator.Check(false, "Err", "ErrMsg")
	validator.Check(false, "Err", "ErrMsg2")

	assert.Equal(t, len(validator.Errors()), 1)
	assert.Equal(t, validator.Errors(), map[string]string{"Err": "ErrMsg"})
}

func TestValidator_Valid(t *testing.T) {
	validator := New()
	assert.True(t, validator.Valid())

	validator.AddError("email", "already exists")
	assert.False(t, validator.Valid())
}

func TestValidator_Match(t *testing.T) {
	tests := []struct {
		value    string
		rx       string
		expected bool
	}{
		{"hello", "hello", true},       // exact
		{"Hello", "hello", false},      // case sensitive
		{"hello world", "hello", true}, // substring
		{"", "hello", false},           // empty string
		{"123", "\\d+", true},          // digits
		{"abc", "\\d+", false},         // non-digits
	}

	for _, tc := range tests {
		actual := Matches(tc.value, regexp.MustCompile(tc.rx))
		assert.Equal(t, tc.expected, actual)
	}
}
