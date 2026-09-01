package schema_validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindConflictingCIDRs(t *testing.T) {
	tests := []struct {
		name           string
		input          []string
		expectedErrors []string
	}{
		{
			name:  "no cidrs",
			input: []string{},
		},
		{
			name:  "single distinct subnet",
			input: []string{"10.0.0.0/8"},
		},
		{
			name:  "distinct subnets",
			input: []string{"10.0.0.0/8", "192.168.0.0/16", "2001:db8::/32"},
		},
		{
			name:  "same subnet with different host bits",
			input: []string{"192.168.10.1/24", "192.168.10.2/24"},
			expectedErrors: []string{
				`CIDRs "192.168.10.1/24", "192.168.10.2/24" belong to the same subnet and would all be stored as "192.168.10.0/24"`,
			},
		},
		{
			name:  "host bit and already-normalized entry",
			input: []string{"192.168.10.1/24", "192.168.10.0/24"},
			expectedErrors: []string{
				`CIDRs "192.168.10.0/24", "192.168.10.1/24" belong to the same subnet and would all be stored as "192.168.10.0/24"`,
			},
		},
		{
			name:  "multiple conflicting groups",
			input: []string{"10.0.0.1/8", "10.1.0.1/8", "192.168.1.1/24", "192.168.1.2/24"},
			expectedErrors: []string{
				`CIDRs "10.0.0.1/8", "10.1.0.1/8" belong to the same subnet and would all be stored as "10.0.0.0/8"`,
				`CIDRs "192.168.1.1/24", "192.168.1.2/24" belong to the same subnet and would all be stored as "192.168.1.0/24"`,
			},
		},
		{
			name:  "invalid cidr is skipped",
			input: []string{"not-a-cidr", "10.0.0.0/8"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findConflictingCIDRs(tt.input)
			assert.Equal(t, tt.expectedErrors, got)
		})
	}
}
