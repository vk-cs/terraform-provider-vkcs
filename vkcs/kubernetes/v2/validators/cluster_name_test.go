package validators

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClusterName(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError error
	}{
		// Valid names
		{
			name:  "valid simple name",
			input: "cluster1",
		},
		{
			name:  "valid name with hyphens",
			input: "cluster-1",
		},
		{
			name:  "valid minimal length (3 chars)",
			input: "abc",
		},
		{
			name:  "valid maximal length (25 chars)",
			input: "a12345678901234567890123",
		},
		{
			name:  "valid name with numbers only",
			input: "12345",
		},

		// Invalid names
		{
			name:          "too short (2 chars)",
			input:         "ab",
			expectedError: errors.New("name length must be no less than 3 and no more than 25 characters, got 2"),
		},
		{
			name:          "too long (26 chars)",
			input:         "a1234567890123456789012343",
			expectedError: errors.New("name length must be no less than 3 and no more than 25 characters, got 26"),
		},

		// Invalid names (spec symbols)
		{
			name:          "contains uppercase letters",
			input:         "Cluster",
			expectedError: errors.New("name must consist of lower case latin alphanumeric characters, '-', must start and end with latin alphanumeric character, and must not contain consecutive hyphens, got invalid characters: [C]"),
		},
		{
			name:          "contains special characters",
			input:         "cluster@1",
			expectedError: errors.New("name must consist of lower case latin alphanumeric characters, '-', must start and end with latin alphanumeric character, and must not contain consecutive hyphens, got invalid characters: [@]"),
		},
		{
			name:          "contains spaces",
			input:         "cluster 1",
			expectedError: errors.New("name must consist of lower case latin alphanumeric characters, '-', must start and end with latin alphanumeric character, and must not contain consecutive hyphens, got invalid characters: [ ]"),
		},
		{
			name:          "starts with hyphen",
			input:         "-cluster",
			expectedError: errors.New("name must consist of lower case latin alphanumeric characters, '-', must start and end with latin alphanumeric character, and must not contain consecutive hyphens"),
		},
		{
			name:          "ends with hyphen",
			input:         "cluster-",
			expectedError: errors.New("name must consist of lower case latin alphanumeric characters, '-', must start and end with latin alphanumeric character, and must not contain consecutive hyphens"),
		},
		{
			name:          "contains consecutive hyphens",
			input:         "cluster--1",
			expectedError: errors.New("consecutive hyphens (--) are not allowed"),
		},

		// Special cases
		{
			name:          "empty string",
			input:         "",
			expectedError: errors.New("name length must be no less than 3 and no more than 25 characters, got 0"),
		},
		{
			name:          "only hyphens",
			input:         "---",
			expectedError: errors.New("consecutive hyphens (--) are not allowed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := isValidClusterName(tt.input)

			if tt.expectedError != nil {
				if assert.Error(t, err) {
					assert.Equal(t, tt.expectedError, err)
				}
				return
			}

			assert.NoError(t, err)
		})
	}
}
