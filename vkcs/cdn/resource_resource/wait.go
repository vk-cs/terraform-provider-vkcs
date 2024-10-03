package resource_resource

import (
	"context"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/cdn/v1/resources"
)

const ResourceReadyTimeout = 10 * time.Minute

const resourceStatePollInterval = 10 * time.Second

func WaitForResourceReady(ctx context.Context, client *gophercloud.ServiceClient, projectID string, resourceID int, timeout time.Duration) diag.Diagnostics {
	var diags diag.Diagnostics

	stateConf := &retry.StateChangeConf{
		Pending:      []string{string(ResourceStatusProcessed)},
		Target:       []string{string(ResourceStatusActive), string(ResourceStatusSuspended)},
		Refresh:      resourceStateRefreshFunc(client, projectID, resourceID),
		Timeout:      timeout,
		PollInterval: resourceStatePollInterval,
	}

	tflog.Trace(ctx, "Waiting for CDN resource to become ready")

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		diags.AddError("Error waiting for CDN resource to become ready", err.Error())
		return diags
	}

	tflog.Trace(ctx, "Waited for CDN resource to become ready")

	return diags
}

func resourceStateRefreshFunc(client *gophercloud.ServiceClient, projectID string, resourceID int) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := resources.Get(client, projectID, resourceID).Extract()
		if err != nil {
			return nil, "", err
		}
		return r, r.Status, nil
	}
}
