package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/v2/validators/helpers"
)

func TestKubernetesLabelsValidator(t *testing.T) {
	tests := []struct {
		name        string
		labels      map[string]string
		shouldError bool
	}{
		// valid labels
		{
			name: "valid simple labels",
			labels: map[string]string{
				"app":     "myapp",
				"version": "v1",
			},
			shouldError: false,
		},
		{
			name: "valid labels with prefix",
			labels: map[string]string{
				"app.kubernetes.io/name":    "myapp",
				"app.kubernetes.io/version": "v1",
			},
			shouldError: false,
		},
		{
			name: "valid labels with dots and hyphens",
			labels: map[string]string{
				"my-app.name": "test-value",
				"env":         "production",
			},
			shouldError: false,
		},
		{
			name: "valid empty value",
			labels: map[string]string{
				"app": "",
			},
			shouldError: false,
		},
		{
			name: "valid single label",
			labels: map[string]string{
				"app": "myapp",
			},
			shouldError: false,
		},
		{
			name: "valid label with underscore",
			labels: map[string]string{
				"app_name": "my_app",
			},
			shouldError: false,
		},
		{
			name: "valid label with uppercase in key",
			labels: map[string]string{
				"App": "myapp",
			},
			shouldError: false,
		},
		{
			name: "valid label with uppercase in value",
			labels: map[string]string{
				"app": "MyApp",
			},
			shouldError: false,
		},
		{
			name: "valid label with consecutive dots",
			labels: map[string]string{
				"app": "my..value",
			},
			shouldError: false,
		},
		{
			name: "valid label with consecutive hyphens",
			labels: map[string]string{
				"app": "my--value",
			},
			shouldError: false,
		},
		{
			name: "valid label with consecutive underscores",
			labels: map[string]string{
				"app": "my__value",
			},
			shouldError: false,
		},
		{
			name: "valid key with consecutive dots",
			labels: map[string]string{
				"my..app": "test",
			},
			shouldError: false,
		},
		{
			name: "valid key with consecutive hyphens",
			labels: map[string]string{
				"my--app": "test",
			},
			shouldError: false,
		},
		{
			name: "valid key with consecutive underscores",
			labels: map[string]string{
				"my__app": "test",
			},
			shouldError: false,
		},

		// invalid keys
		{
			name: "invalid key with special characters",
			labels: map[string]string{
				"app@name": "myapp",
			},
			shouldError: true,
		},
		{
			name: "invalid key starting with hyphen",
			labels: map[string]string{
				"-app": "myapp",
			},
			shouldError: true,
		},
		{
			name: "invalid key ending with hyphen",
			labels: map[string]string{
				"app-": "myapp",
			},
			shouldError: true,
		},
		{
			name: "invalid key with multiple slashes",
			labels: map[string]string{
				"app/name/version": "myapp",
			},
			shouldError: true,
		},
		{
			name: "invalid key with empty prefix",
			labels: map[string]string{
				"/name": "myapp",
			},
			shouldError: true,
		},
		{
			name: "invalid key with empty name",
			labels: map[string]string{
				"app/": "myapp",
			},
			shouldError: true,
		},
		{
			name: "invalid key too long",
			labels: map[string]string{
				"a234567890123456789012345678901234567890123456789012345678901234": "test",
			},
			shouldError: true,
		},

		// invalid values
		{
			name: "invalid value with special characters",
			labels: map[string]string{
				"app": "my@app",
			},
			shouldError: true,
		},
		{
			name: "invalid value starting with hyphen",
			labels: map[string]string{
				"app": "-myapp",
			},
			shouldError: true,
		},
		{
			name: "invalid value ending with hyphen",
			labels: map[string]string{
				"app": "myapp-",
			},
			shouldError: true,
		},
		{
			name: "invalid value too long",
			labels: map[string]string{
				"app": "thisisareallylonglabelvaluethatexceedsthesixtythreecharacterlimitandshouldfailvalidation",
			},
			shouldError: true,
		},

		// mixed valid and invalid
		{
			name: "mixed valid and invalid keys",
			labels: map[string]string{
				"app":      "myapp",
				"app@name": "myapp", // invalid
				"version":  "v1",
			},
			shouldError: true,
		},
		{
			name: "mixed valid and invalid values",
			labels: map[string]string{
				"app":     "myapp",
				"version": "v1@", // invalid
				"env":     "prod",
			},
			shouldError: true,
		},
		{
			name: "multiple invalid entries",
			labels: map[string]string{
				"-app":    "myapp", // invalid key
				"env":     "-prod", // invalid value
				"version": "v1",
			},
			shouldError: true,
		},

		// edge cases
		{
			name:        "empty map",
			labels:      map[string]string{},
			shouldError: false,
		},
		{
			name: "key at max length",
			labels: map[string]string{
				"a23456789012345678901234567890123456789012345678901234567890123": "test",
			},
			shouldError: false,
		},
		{
			name: "value at max length",
			labels: map[string]string{
				"app": "a23456789012345678901234567890123456789012345678901234567890123",
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// validate each key-value pair
			hasError := false
			for key, value := range tt.labels {
				// validate key
				if err := helpers.IsQualifiedName(key); err != nil {
					hasError = true
					break
				}

				// validate value
				if err := helpers.IsValidLabelValue(value); err != nil {
					hasError = true
					break
				}
			}

			if tt.shouldError {
				assert.True(t, hasError, "expected validation error but got none for labels: %v", tt.labels)
			} else {
				assert.False(t, hasError, "unexpected validation error for labels: %v", tt.labels)
			}
		})
	}
}
