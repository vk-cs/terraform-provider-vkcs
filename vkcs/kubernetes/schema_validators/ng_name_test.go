package schema_validators

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidNodeGroupName(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError error
	}{
		// Valid names
		{
			name:  "valid simple name",
			input: "nodegroup1",
		},
		{
			name:  "valid name with hyphens",
			input: "node-group-1",
		},
		{
			name:  "valid minimal length (3 chars)",
			input: "ng1",
		},
		{
			name:  "valid maximal length (25 chars)",
			input: "a12345678901234567890123",
		},
		{
			name:  "valid name with numbers only",
			input: "12345",
		},
		{
			name:  "valid name 'workers'",
			input: "workers",
		},
		{
			name:  "valid name 'nodes'",
			input: "nodes",
		},

		// Invalid names (blacklisted names)
		{
			name:          "blacklisted name 'master'",
			input:         "master",
			expectedError: errors.New("Node group name can't be a word 'master'"),
		},
		{
			name:          "blacklisted name 'master' with different case",
			input:         "Master",
			expectedError: errors.New("Node group name must consist of lower case Latin alphanumeric characters, '-', must start and end with Latin alphanumeric character, and must not contain consecutive hyphens, got invalid characters: ['M']"),
		},
		{
			name:  "contains 'master' as part of name",
			input: "master-nodes",
		},
		{
			name:          "too short (2 chars)",
			input:         "ng",
			expectedError: errors.New("Node group name length must be no less than 3 and no more than 25 characters, got 2"),
		},
		{
			name:          "too long (26 chars)",
			input:         "a1234567890123456789012345",
			expectedError: errors.New("Node group name length must be no less than 3 and no more than 25 characters, got 26"),
		},
		{
			name:          "contains uppercase letters",
			input:         "NodeGroup",
			expectedError: errors.New("Node group name must consist of lower case Latin alphanumeric characters, '-', must start and end with Latin alphanumeric character, and must not contain consecutive hyphens, got invalid characters: ['N' 'G']"),
		},
		{
			name:          "contains special characters",
			input:         "node@group",
			expectedError: errors.New("Node group name must consist of lower case Latin alphanumeric characters, '-', must start and end with Latin alphanumeric character, and must not contain consecutive hyphens, got invalid characters: ['@']"),
		},
		{
			name:          "contains spaces",
			input:         "node group",
			expectedError: errors.New("Node group name must consist of lower case Latin alphanumeric characters, '-', must start and end with Latin alphanumeric character, and must not contain consecutive hyphens, got invalid characters: [' ']"),
		},
		{
			name:          "starts with hyphen",
			input:         "-nodegroup",
			expectedError: errors.New("Node group name must consist of lower case Latin alphanumeric characters, '-', must start and end with Latin alphanumeric character, and must not contain consecutive hyphens"),
		},
		{
			name:          "ends with hyphen",
			input:         "nodegroup-",
			expectedError: errors.New("Node group name must consist of lower case Latin alphanumeric characters, '-', must start and end with Latin alphanumeric character, and must not contain consecutive hyphens"),
		},
		{
			name:          "contains consecutive hyphens",
			input:         "node--group",
			expectedError: errors.New("Consecutive hyphens (--) are not allowed"),
		},
		{
			name:          "empty string",
			input:         "",
			expectedError: errors.New("Node group name length must be no less than 3 and no more than 25 characters, got 0"),
		},
		{
			name:          "only hyphens",
			input:         "---",
			expectedError: errors.New("Consecutive hyphens (--) are not allowed"),
		},
		{
			name:          "whitespace only",
			input:         "   ",
			expectedError: errors.New("Node group name length must be no less than 3 and no more than 25 characters, got 0"),
		},
		{
			name:  "name with leading/trailing spaces",
			input: "  nodegroup  ",
		},
		{
			name:  "name 'masters' (contains 'master' but not equal)",
			input: "masters",
		},
		{
			name:  "name 'mastery' (contains 'master' but not equal)",
			input: "mastery",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := isValidNodeGroupName(tt.input)

			if tt.expectedError != nil {
				if assert.Error(t, err) {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				}
				return
			}

			assert.NoError(t, err)
		})
	}
}
