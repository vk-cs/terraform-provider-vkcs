package compute_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/compute"
)

func TestComputeInterfaceAttachParseID(t *testing.T) {
	id := "foo/bar"

	expectedInstanceID := "foo"
	expectedAttachmentID := "bar"

	actualInstanceID, actualAttachmentID, err := compute.ComputeInterfaceAttachParseID(id)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedInstanceID, actualInstanceID)
	assert.Equal(t, expectedAttachmentID, actualAttachmentID)
}
