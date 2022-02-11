package vkcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeInterfaceAttachParseID(t *testing.T) {
	id := "foo/bar"

	expectedInstanceID := "foo"
	expectedAttachmentID := "bar"

	actualInstanceID, actualAttachmentID, err := computeInterfaceAttachParseID(id)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedInstanceID, actualInstanceID)
	assert.Equal(t, expectedAttachmentID, actualAttachmentID)
}
