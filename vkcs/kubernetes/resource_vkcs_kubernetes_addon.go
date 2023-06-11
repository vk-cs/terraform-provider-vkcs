package kubernetes

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	v1 "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfraaddons/v1"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfraaddons/v1/addons"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfraaddons/v1/clusteraddons"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

const (
	addonDelay         = 10 * time.Second
	addonMinTimeout    = 10 * time.Second
	addonCreateTimeout = 30 * time.Minute
	addonDeleteTimeout = 30 * time.Minute
)

const (
	addonStatusNew        = "NEW"
	addonStatusInstalling = "INSTALLING"
	addonStatusInstalled  = "INSTALLED"
	addonStatusReplaced   = "REPLACED"
	addonStatusDeleting   = "DELETING"
	addonStatusDeleted    = "DELETED"
	addonStatusFailed     = "FAILED"
)

var _ resource.Resource = &AddonResource{}
var _ resource.ResourceWithConfigure = &AddonResource{}
var _ resource.ResourceWithImportState = &AddonResource{}

func NewAddonResource() resource.Resource {
	return &AddonResource{}
}

type AddonResource struct {
	config clients.Config
}

type AddonResourceModel struct {
	ID                  types.String   `tfsdk:"id"`
	Region              types.String   `tfsdk:"region"`
	ClusterID           types.String   `tfsdk:"cluster_id"`
	AddonID             types.String   `tfsdk:"addon_id"`
	Namespace           types.String   `tfsdk:"namespace"`
	Name                types.String   `tfsdk:"name"`
	ConfigurationValues types.String   `tfsdk:"configuration_values"`
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
}

func (r *AddonResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_kubernetes_addon"
}

func (r *AddonResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource",
			},

			"region": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The region in which to obtain the Container Infra Addons client. If omitted, the `region` argument of the provider is used. Changing this creates a new addon.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"cluster_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the kubernetes cluster. Changing this creates a new addon.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"addon_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the addon. Changing this creates a new addon.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"namespace": schema.StringAttribute{
				Required:    true,
				Description: "The namespace name where the addon will be installed.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the application. Changing this creates a new addon.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"configuration_values": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Configuration code for the addon. Changing this creates a new addon.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Delete: true,
			}),
		},

		Description: "Provides a kubernetes cluster addon resource. This can be used to create, modify and delete kubernetes cluster addons.",
	}
}

func (r *AddonResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *AddonResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *AddonResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, addonCreateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	client, err := r.config.ContainerInfraAddonsV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Kubernetes Addons API client", err.Error())
		return
	}

	addonID, clusterID := data.AddonID.ValueString(), data.ClusterID.ValueString()
	ctx = tflog.SetField(ctx, "cluster_id", clusterID)
	ctx = tflog.SetField(ctx, "addon_id", addonID)

	name := data.Name.ValueString()
	if name == "" {
		tflog.Debug(ctx, "Calling Addons API to resolve addon name by its ID")

		availableAddon, err := addons.GetAvailableAddon(client, clusterID, addonID).Extract()
		if err != nil {
			resp.Diagnostics.AddError("Error calling Kubernetes Addons API", err.Error())
			return
		}

		name = availableAddon.Name

		tflog.Debug(ctx, "Called Addons API to resolve addon name by its ID", map[string]interface{}{"available_addon": fmt.Sprintf("%#v", availableAddon)})
	}

	createOpts := addons.InstallAddonToClusterOpts{
		Values: data.ConfigurationValues.ValueString(),
		Payload: v1.Payload{
			Namespace: data.Namespace.ValueString(),
			Name:      name,
		},
	}

	tflog.Debug(ctx, "Calling Addons API to install the addon to the cluster")

	clusterAddon, err := addons.InstallAddonToCluster(client, addonID, clusterID, &createOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling Kubernetes Addons API", err.Error())
		return
	}

	tflog.Debug(ctx, "Called Addons API to install the addon to the cluster", map[string]interface{}{"cluster_addon": fmt.Sprintf("%#v", clusterAddon)})

	id := clusterAddon.ID
	resp.State.SetAttribute(ctx, path.Root("id"), id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{addonStatusNew, addonStatusInstalling},
		Target:     []string{addonStatusInstalled},
		Refresh:    addonStateRefreshFunc(client, id),
		Timeout:    timeout,
		Delay:      addonDelay,
		MinTimeout: addonMinTimeout,
	}

	tflog.Debug(ctx, "Waiting for the addon to be installed", map[string]interface{}{"timeout": timeout})

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for the addon to become ready", err.Error())
		return
	}

	data.ID = types.StringValue(id)
	data.Region = types.StringValue(region)
	data.AddonID = types.StringValue(clusterAddon.Addon.ID)
	data.Namespace = types.StringValue(clusterAddon.Payload.Namespace)
	data.Name = types.StringValue(clusterAddon.Payload.Name)
	data.ConfigurationValues = types.StringValue(clusterAddon.UserChartValues)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AddonResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AddonResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	client, err := r.config.ContainerInfraAddonsV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Kubernetes Addons API client", err.Error())
		return
	}

	id := data.ID.ValueString()
	ctx = tflog.SetField(ctx, "cluster_addon_id", id)

	tflog.Debug(ctx, "Calling Addons API to retrieve cluster addon")

	clusterAddon, err := clusteraddons.Get(client, id).Extract()
	if errutil.Is(err, 404) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error calling Kubernetes Addons API", err.Error())
		return
	}

	tflog.Debug(ctx, "Called Addons API to retrieve cluster addon", map[string]interface{}{"cluster_addon": fmt.Sprintf("%#v", clusterAddon)})

	data.ID = types.StringValue(clusterAddon.ID)
	data.Region = types.StringValue(region)
	data.AddonID = types.StringValue(clusterAddon.Addon.ID)
	data.Namespace = types.StringValue(clusterAddon.Payload.Namespace)
	data.Name = types.StringValue(clusterAddon.Payload.Name)
	data.ConfigurationValues = types.StringValue(clusterAddon.UserChartValues)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AddonResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *AddonResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("Unable to update the addon",
		"Not implemented. Please report this issue to the provider developers.")
}

func (r *AddonResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *AddonResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, addonDeleteTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	client, err := r.config.ContainerInfraAddonsV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Kubernetes Addons API client", err.Error())
		return
	}

	id := data.ID.ValueString()
	ctx = tflog.SetField(ctx, "cluster_addon_id", id)

	tflog.Debug(ctx, "Calling Addons API to check if the addon has been already deleted")

	clusterAddon, err := clusteraddons.Get(client, id).Extract()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error calling Kubernetes Addons API", err.Error())
		return
	}
	if clusterAddon.Status == addonStatusDeleted {
		return
	}

	tflog.Debug(ctx, "The addon is still present. Calling Addons API to delete it")

	err = clusteraddons.Delete(client, id).ExtractErr()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error calling Kubernetes Addons API", err.Error())
		return
	}

	tflog.Debug(ctx, "Called Addons API to delete the addon")

	stateConf := &retry.StateChangeConf{
		Pending:    []string{addonStatusDeleting},
		Target:     []string{addonStatusDeleted},
		Refresh:    addonStateRefreshFunc(client, id),
		Timeout:    timeout,
		Delay:      addonDelay,
		MinTimeout: addonMinTimeout,
	}

	tflog.Debug(ctx, "Waiting for the addon to be deleted", map[string]interface{}{"timeout": timeout})

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for the addon to become ready", err.Error())
		return
	}
}

func (r *AddonResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: cluster_id/cluster_addon_id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cluster_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}
