package lb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandLBListenerHeadersMap(t *testing.T) {
	raw := map[string]interface{}{
		"header0": "val0",
		"header1": "val1",
	}

	expected := map[string]string{
		"header0": "val0",
		"header1": "val1",
	}

	actual, err := expandLBListenerHeadersMap(raw)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestExpandLBListenerHeadersMap_err(t *testing.T) {
	raw := map[string]interface{}{
		"header0": "val0",
		"header1": 1,
	}

	actual, err := expandLBListenerHeadersMap(raw)

	assert.Error(t, err)
	assert.Empty(t, actual)
}
