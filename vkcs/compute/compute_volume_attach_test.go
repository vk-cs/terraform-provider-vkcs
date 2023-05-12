package compute_test

import (
	"testing"

	"github.com/vk-cs/terraform-provider-vkcs/vkcs/compute"
)

func TestComputeVolumeAttachV2ParseID(t *testing.T) {
	id := "foo/bar"

	expectedInstanceID := "foo"
	expectedAttachmentID := "bar"

	actualInstanceID, actualAttachmentID, err := compute.ComputeVolumeAttachParseID(id)

	if err != nil {
		t.Fatal(err)
	}

	if expectedInstanceID != actualInstanceID {
		t.Fatalf("Instance IDs differ. Want %s, but got %s", expectedInstanceID, actualInstanceID)
	}

	if expectedAttachmentID != actualAttachmentID {
		t.Fatalf("Attachment IDs differ. Want %s, but got %s", expectedAttachmentID, actualAttachmentID)
	}
}
