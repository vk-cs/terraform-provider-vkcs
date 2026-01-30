package validators

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidClusterVersion(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError error
	}{
		// Valid versions
		{
			name:  "valid version v1.34.2",
			input: "v1.34.2",
		},
		{
			name:  "valid version v1.0.0",
			input: "v1.0.0",
		},
		{
			name:  "valid version v1.34.0",
			input: "v1.34.0",
		},
		{
			name:  "valid version v2.0.0",
			input: "v2.0.0",
		},
		{
			name:  "valid version v1.100.200",
			input: "v1.100.200",
		},
		{
			name:  "valid version with patch version 0",
			input: "v1.34.0",
		},
		{
			name:  "valid version with minor version 0",
			input: "v1.0.15",
		},

		// Invalid versions (format)
		{
			name:          "missing v prefix",
			input:         "1.34.2",
			expectedError: errors.New("version '1.34.2' does not match expected format vX.XX.XX where X is a number (e.g., v1.34.2)"),
		},
		{
			name:          "uppercase V",
			input:         "V1.34.2",
			expectedError: errors.New("version 'V1.34.2' does not match expected format vX.XX.XX where X is a number (e.g., v1.34.2)"),
		},
		{
			name:          "too many components",
			input:         "v1.34.2.1",
			expectedError: errors.New("version 'v1.34.2.1' does not match expected format vX.XX.XX where X is a number (e.g., v1.34.2)"),
		},
		{
			name:          "too few components",
			input:         "v1.34",
			expectedError: errors.New("version 'v1.34' does not match expected format vX.XX.XX where X is a number (e.g., v1.34.2)"),
		},
		{
			name:          "single component",
			input:         "v1",
			expectedError: errors.New("version 'v1' does not match expected format vX.XX.XX where X is a number (e.g., v1.34.2)"),
		},
		{
			name:          "empty string",
			input:         "",
			expectedError: errors.New("version '' does not match expected format vX.XX.XX where X is a number (e.g., v1.34.2)"),
		},
		{
			name:          "with build metadata",
			input:         "v1.34.2+build",
			expectedError: errors.New("version 'v1.34.2+build' does not match expected format vX.XX.XX where X is a number (e.g., v1.34.2)"),
		},
		{
			name:          "with prerelease",
			input:         "v1.34.2-alpha",
			expectedError: errors.New("version 'v1.34.2-alpha' does not match expected format vX.XX.XX where X is a number (e.g., v1.34.2)"),
		},

		// Invalid versions (non-numeric components)
		{
			name:          "non-numeric major",
			input:         "va.34.2",
			expectedError: errors.New("version 'va.34.2' does not match expected format vX.XX.XX where X is a number (e.g., v1.34.2)"),
		},
		{
			name:          "non-numeric minor",
			input:         "v1.xx.2",
			expectedError: errors.New("version 'v1.xx.2' does not match expected format vX.XX.XX where X is a number (e.g., v1.34.2)"),
		},
		{
			name:          "non-numeric patch",
			input:         "v1.34.x",
			expectedError: errors.New("version 'v1.34.x' does not match expected format vX.XX.XX where X is a number (e.g., v1.34.2)"),
		},
		{
			name:          "contains letters in components",
			input:         "v1.3a.2",
			expectedError: errors.New("version 'v1.3a.2' does not match expected format vX.XX.XX where X is a number (e.g., v1.34.2)"),
		},

		// Invalid versions (negative components) - ИЗМЕНЕНО: теперь должны возвращаться конкретные ошибки
		{
			name:          "negative minor version",
			input:         "v1.-1.2",
			expectedError: errors.New("minor version of cluster version cannot be negative, got -1"),
		},
		{
			name:          "negative patch version",
			input:         "v1.34.-2",
			expectedError: errors.New("patch version of cluster version cannot be negative, got -2"),
		},
		{
			name:          "negative major version",
			input:         "v-1.34.2",
			expectedError: errors.New("major version of cluster version must be at least 1, got -1"),
		},
		{
			name:          "negative minor version with leading zeros",
			input:         "v1.-02.3",
			expectedError: errors.New("minor version of cluster version cannot be negative, got -2"),
		},
		{
			name:          "negative patch version with leading zeros",
			input:         "v1.34.-02",
			expectedError: errors.New("patch version of cluster version cannot be negative, got -2"),
		},
		{
			name:          "negative major version with leading zeros",
			input:         "v-01.2.3",
			expectedError: errors.New("major version of cluster version must be at least 1, got -1"),
		},
		{
			name:          "all components negative",
			input:         "v-1.-2.-3",
			expectedError: errors.New("major version of cluster version must be at least 1, got -1"),
		},

		// Invalid versions (logical components)
		{
			name:          "major version 0",
			input:         "v0.34.2",
			expectedError: errors.New("major version of cluster version must be at least 1, got 0"),
		},
		{
			name:          "major version 0 with leading zeros",
			input:         "v00.34.2",
			expectedError: errors.New("major version of cluster version must be at least 1, got 0"),
		},

		// Special cases
		{
			name:          "version with leading zeros",
			input:         "v01.34.2",
			expectedError: errors.New("version 'v01.34.2' contains leading zeros in version components; version components must not have leading zeros"),
		},
		{
			name:          "version with spaces",
			input:         "v1.34. 2",
			expectedError: errors.New("version 'v1.34. 2' does not match expected format vX.XX.XX where X is a number (e.g., v1.34.2)"),
		},
		{
			name:          "version with extra dots",
			input:         "v1..34.2",
			expectedError: errors.New("version 'v1..34.2' does not match expected format vX.XX.XX where X is a number (e.g., v1.34.2)"),
		},
		{
			name:          "version with special characters",
			input:         "v1.34_2",
			expectedError: errors.New("version 'v1.34_2' does not match expected format vX.XX.XX where X is a number (e.g., v1.34.2)"),
		},
		{
			name:  "valid version with multiple digits in all components",
			input: "v123.456.789",
		},
		{
			name:          "version with leading zero in minor",
			input:         "v1.034.2",
			expectedError: errors.New("version 'v1.034.2' contains leading zeros in version components; version components must not have leading zeros"),
		},
		{
			name:          "version with leading zero in patch",
			input:         "v1.34.02",
			expectedError: errors.New("version 'v1.34.02' contains leading zeros in version components; version components must not have leading zeros"),
		},
		{
			name:          "version with multiple leading zeros",
			input:         "v001.034.002",
			expectedError: errors.New("version 'v001.034.002' contains leading zeros in version components; version components must not have leading zeros"),
		},
		// Edge cases
		{
			name:          "version with minus in the middle",
			input:         "v1.-.2",
			expectedError: errors.New("version 'v1.-.2' does not match expected format vX.XX.XX where X is a number (e.g., v1.34.2)"),
		},
		{
			name:          "version with only minus sign",
			input:         "v1.-.2",
			expectedError: errors.New("version 'v1.-.2' does not match expected format vX.XX.XX where X is a number (e.g., v1.34.2)"),
		},
		{
			name:          "version with minus at the end",
			input:         "v1.2.-",
			expectedError: errors.New("version 'v1.2.-' does not match expected format vX.XX.XX where X is a number (e.g., v1.34.2)"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := isValidClusterVersion(tt.input)

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
