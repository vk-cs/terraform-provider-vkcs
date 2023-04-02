package vkcs

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatabaseCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	if diff.Id() != "" && diff.HasChange("cloud_monitoring_enabled") {
		t, exists := diff.GetOk("datastore.0.type")
		if !exists {
			return errors.New("datastore.0.type is not found")
		}
		if exists && isOperationNotSupported(t.(string), Redis, MongoDB) {
			return diff.ForceNew("cloud_monitoring_enabled")
		}
	}
	return nil
}
