package schema_validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInsecureRegistries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value []string
		err   bool
	}{
		{
			name:  "domain",
			value: []string{"myregistry.com"},
			err:   false,
		},
		{
			name:  "domain without dot",
			value: []string{"myregistry"},
			err:   true,
		},
		{
			name:  "ip",
			value: []string{"192.255.0.12"},
			err:   false,
		},
		{
			name:  "invalid ip",
			value: []string{"192.256.0.12"},
			err:   true,
		},
		{
			name:  "domain with port",
			value: []string{"myregistry.com:5000"},
			err:   false,
		},
		{
			name:  "domain with incorrect port",
			value: []string{"myregistry.com:65536"},
			err:   true,
		},
		{
			name:  "domain with empty path",
			value: []string{"example.com/"},
			err:   true,
		},
		{
			name:  "domain with consecutive slashes",
			value: []string{"example.com/myorg//departament"},
			err:   true,
		},
		{
			name:  "path",
			value: []string{"myregistry.com/myorg/department/myteam"},
			err:   false,
		},
		{
			name:  "domain with port and path",
			value: []string{"myregistry.com:5000/myorg/department/myteam"},
			err:   false,
		},
		{
			name:  "path with space",
			value: []string{"myregistry.com/my org/department/myteam"},
			err:   true,
		},
		{
			name:  "path with quote",
			value: []string{"myregistry.com/my\"org/department/myteam"},
			err:   true,
		},
		{
			name:  "image",
			value: []string{"myregistry:latest"},
			err:   true,
		},
		{
			name:  "image with tag",
			value: []string{"myregistry.com/myorg/myimage:1.0"},
			err:   true,
		},
		{
			name:  "image with digest",
			value: []string{"myregistry.com/myorg/image@sha256:xyz"},
			err:   true,
		},
		{
			name:  "wildcard",
			value: []string{"*.example.com"},
			err:   false,
		},
		{
			name:  "wildcard in the middle",
			value: []string{"example.*.com"},
			err:   true,
		},
		{
			name:  "wildcard with path",
			value: []string{"*.example.com/foo"},
			err:   true,
		},
		{
			name:  "wildcard with port",
			value: []string{"*.example.com:5000"},
			err:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			for _, insecReg := range tt.value {
				err := isValidInsecureRegistry(insecReg)
				if tt.err {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}
