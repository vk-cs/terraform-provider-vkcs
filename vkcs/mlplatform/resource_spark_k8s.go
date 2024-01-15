package mlplatform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/framework/planmodifiers"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/mlplatform/v1/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

// Ensure the implementation satisfies the desired interfaces.
var _ resource.Resource = &SparkK8SResource{}

func NewSparkK8SResource() resource.Resource {
	return &SparkK8SResource{}
}

type SparkK8SResource struct {
	config clients.Config
}

func (r *SparkK8SResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_mlplatform_spark_k8s"
}

type SparkK8SResourceModel struct {
	ID                   types.String              `tfsdk:"id"`
	Name                 types.String              `tfsdk:"name"`
	AvailabilityZone     types.String              `tfsdk:"availability_zone"`
	NetworkID            types.String              `tfsdk:"network_id"`
	SubnetID             types.String              `tfsdk:"subnet_id"`
	NodeGroups           []*SparkK8SNodeGroupModel `tfsdk:"node_groups"`
	ClusterMode          types.String              `tfsdk:"cluster_mode"`
	RegistryID           types.String              `tfsdk:"registry_id"`
	IPPool               types.String              `tfsdk:"ip_pool"`
	SparkConfiguration   types.String              `tfsdk:"spark_configuration"`
	EnvironmentVariables types.String              `tfsdk:"environment_variables"`

	S3BucketName            types.String `tfsdk:"s3_bucket_name"`
	HistoryServerURL        types.String `tfsdk:"history_server_url"`
	ControlInstanceID       types.String `tfsdk:"control_instance_id"`
	InactiveMin             types.Int64  `tfsdk:"inactive_min"`
	SuspendAfterInactiveMin types.Int64  `tfsdk:"suspend_after_inactive_min"`
	DeleteAfterInactiveMin  types.Int64  `tfsdk:"delete_after_inactive_min"`

	Region   types.String   `tfsdk:"region"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type SparkK8SNodeGroupModel struct {
	NodeCount          types.Int64  `tfsdk:"node_count"`
	FlavorID           types.String `tfsdk:"flavor_id"`
	AutoscalingEnabled types.Bool   `tfsdk:"autoscaling_enabled"`
	MinNodes           types.Int64  `tfsdk:"min_nodes"`
	MaxNodes           types.Int64  `tfsdk:"max_nodes"`
}

func (r *SparkK8SResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"name": schema.StringAttribute{
				Required:    true,
				Description: "Cluster name. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"availability_zone": schema.StringAttribute{
				Required:    true,
				Description: "The availability zone in which to create the resource. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"network_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the network. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"subnet_id": schema.StringAttribute{
				Optional:    true,
				Description: "ID of the subnet. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"node_groups": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"node_count": schema.Int64Attribute{
							Optional:    true,
							Description: "Count of nodes in node group",
						},

						"flavor_id": schema.StringAttribute{
							Required:    true,
							Description: "ID of the flavor to be used in nodes",
						},

						"autoscaling_enabled": schema.BoolAttribute{
							Required:    true,
							Description: "Enables autoscaling for node group",
						},

						"min_nodes": schema.Int64Attribute{
							Optional:    true,
							Description: "Minimum count of nodes in node group. It is used only when autoscaling is enabled",
						},

						"max_nodes": schema.Int64Attribute{
							Optional:    true,
							Description: "Maximum number of nodes in node group. It is used only when autoscaling is enabled",
						},
					},
				},
				Required:    true,
				Description: "Cluster's node groups configuration",
			},

			"cluster_mode": schema.StringAttribute{
				Required:    true,
				Description: "Cluster Mode. Should be `DEV` or `PROD`. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("DEV", "PROD"),
				},
			},

			"registry_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the K8S registry to use with cluster. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"ip_pool": schema.StringAttribute{
				Required:    true,
				Description: "ID of the ip pool. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"spark_configuration": schema.StringAttribute{
				Optional:    true,
				Description: "Spark configuration. Read more about this parameter here: https://cloud.vk.com/docs/en/ml/spark-to-k8s/instructions/create. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"environment_variables": schema.StringAttribute{
				Optional:    true,
				Description: "Environment variables. Read more about this parameter here: https://cloud.vk.com/docs/en/ml/spark-to-k8s/instructions/create. Changing this creates a new resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"suspend_after_inactive_min": schema.Int64Attribute{
				Optional:    true,
				Description: "Timeout of cluster inactivity before suspending, in minutes. Changing this creates a new resource",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},

			"delete_after_inactive_min": schema.Int64Attribute{
				Optional:    true,
				Description: "Timeout of cluster inactivity before deletion, in minutes. Changing this creates a new resource",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},

			"s3_bucket_name": schema.StringAttribute{
				Computed:    true,
				Description: "S3 bucket name",
			},

			"history_server_url": schema.StringAttribute{
				Computed:    true,
				Description: "URL of the history server",
			},

			"control_instance_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the control instance",
			},

			"inactive_min": schema.Int64Attribute{
				Computed:    true,
				Description: "Current time of cluster inactivity, in minutes",
			},

			"region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The `region` in which ML Platform client is obtained, defaults to the provider's `region`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIf(planmodifiers.GetRegionPlanModifier(resp),
						"require replacement if configuration value changes", "require replacement if configuration value changes"),
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
		Description: "Manages a ML Platform Spark K8S cluster resource.",
	}
}

func (r *SparkK8SResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *SparkK8SResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SparkK8SResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, clusterCreateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	mlPlatformClient, err := r.config.MLPlatformV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS ML Platform client", err.Error())
		return
	}

	nodeGroups := expandNodeGroupOpts(data.NodeGroups)

	sparkK8SCreateOpts := clusters.CreateOpts{
		Name:             data.Name.ValueString(),
		AvailabilityZone: data.AvailabilityZone.ValueString(),
		NetworkID:        data.NetworkID.ValueString(),
		SubnetID:         data.SubnetID.ValueString(),
		NodeGroups:       nodeGroups,
		ClusterMode:      data.ClusterMode.ValueString(),
		Registry: clusters.RegistryCreateOpts{
			ExistingRegistryID: data.RegistryID.ValueString(),
		},
		IPPool:               data.IPPool.ValueString(),
		DeleteAfterDelay:     int(data.DeleteAfterInactiveMin.ValueInt64()) * 60,
		SuspendAfterDelay:    int(data.SuspendAfterInactiveMin.ValueInt64()) * 60,
		SparkConfiguration:   data.SparkConfiguration.ValueString(),
		EnvironmentVariables: data.EnvironmentVariables.ValueString(),
	}

	tflog.Debug(ctx, "Creating Spark K8S Cluster", map[string]interface{}{"options": fmt.Sprintf("%#v", sparkK8SCreateOpts)})

	sparkK8SResp, err := clusters.Create(mlPlatformClient, &sparkK8SCreateOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_mlplatform_spark_k8s", err.Error())
		return
	}
	resp.State.SetAttribute(ctx, path.Root("id"), sparkK8SResp.ID)

	sparkK8SStateConf := &retry.StateChangeConf{
		Pending:    []string{clusterStatusCreating, clusterStatusInstallSparkOperator, clusterStatusProvisioning, clusterStatusStarting, clusterStatusCreatingRegistry},
		Target:     []string{clusterStatusRunning},
		Refresh:    clusterStateRefreshFunc(mlPlatformClient, sparkK8SResp.ID),
		Timeout:    timeout,
		Delay:      clusterDelay,
		MinTimeout: clusterMinTimeout,
	}

	tflog.Debug(ctx, "Waiting for the Spark K8S cluster to be created", map[string]interface{}{"timeout": timeout})

	_, err = sparkK8SStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for the cluster to become ready", err.Error())
		return
	}

	sparkK8SResp, err = clusters.Get(mlPlatformClient, sparkK8SResp.ID).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vkcs_mlplatform_spark_k8s", err.Error())
		return
	}

	data.ID = types.StringValue(sparkK8SResp.ID)
	data.Region = types.StringValue(region)

	data.Name = types.StringValue(sparkK8SResp.Name)
	data.AvailabilityZone = types.StringValue(sparkK8SResp.Info.AvailabilityZone)
	data.NetworkID = types.StringValue(sparkK8SResp.Info.NetworkID)
	data.SubnetID = types.StringValue(sparkK8SResp.Info.SubnetID)
	data.ClusterMode = types.StringValue(sparkK8SResp.Info.ClusterMode)
	data.RegistryID = types.StringValue(sparkK8SResp.DockerRegistryID)
	data.IPPool = types.StringValue(sparkK8SResp.Info.IPPool)

	data.S3BucketName = types.StringValue(sparkK8SResp.S3BucketName)
	data.HistoryServerURL = types.StringValue(sparkK8SResp.HistoryServerURL)
	data.ControlInstanceID = types.StringValue(sparkK8SResp.ControlInstanceID)
	data.InactiveMin = types.Int64Value(int64(sparkK8SResp.InactiveMin))
	data.SuspendAfterInactiveMin = types.Int64Value(int64(sparkK8SResp.SuspendAfterInactiveMin))
	data.DeleteAfterInactiveMin = types.Int64Value(int64(sparkK8SResp.DeleteAfterInactiveMin))
	data.SparkConfiguration = types.StringValue(sparkK8SResp.Info.SparkConfiguration)
	data.EnvironmentVariables = types.StringValue(sparkK8SResp.Info.EnvironmentVariables)

	data.NodeGroups = flattenNodeGroupOpts(sparkK8SResp.Info.NodeGroups)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *SparkK8SResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SparkK8SResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	mlPlatformClient, err := r.config.MLPlatformV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS ML Platform client", err.Error())
		return
	}

	sparkK8SID := data.ID.ValueString()

	sparkK8SResp, err := clusters.Get(mlPlatformClient, sparkK8SID).Extract()
	if err != nil {
		checkDeleted := util.CheckDeletedResource(ctx, resp, err)
		if checkDeleted != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_mlplatform_spark_k8s", checkDeleted.Error())
		}
		return
	}
	data.Region = types.StringValue(region)

	data.Name = types.StringValue(sparkK8SResp.Name)
	data.AvailabilityZone = types.StringValue(sparkK8SResp.Info.AvailabilityZone)
	data.NetworkID = types.StringValue(sparkK8SResp.Info.NetworkID)
	data.SubnetID = types.StringValue(sparkK8SResp.Info.SubnetID)
	data.ClusterMode = types.StringValue(sparkK8SResp.Info.ClusterMode)
	data.RegistryID = types.StringValue(sparkK8SResp.DockerRegistryID)
	data.IPPool = types.StringValue(sparkK8SResp.Info.IPPool)

	data.S3BucketName = types.StringValue(sparkK8SResp.S3BucketName)
	data.HistoryServerURL = types.StringValue(sparkK8SResp.HistoryServerURL)
	data.ControlInstanceID = types.StringValue(sparkK8SResp.ControlInstanceID)
	data.InactiveMin = types.Int64Value(int64(sparkK8SResp.InactiveMin))
	data.SuspendAfterInactiveMin = types.Int64Value(int64(sparkK8SResp.SuspendAfterInactiveMin))
	data.DeleteAfterInactiveMin = types.Int64Value(int64(sparkK8SResp.DeleteAfterInactiveMin))
	data.SparkConfiguration = types.StringValue(sparkK8SResp.Info.SparkConfiguration)
	data.EnvironmentVariables = types.StringValue(sparkK8SResp.Info.EnvironmentVariables)

	data.NodeGroups = flattenNodeGroupOpts(sparkK8SResp.Info.NodeGroups)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SparkK8SResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *SparkK8SResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SparkK8SResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, clusterDeleteTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	mlPlatformClient, err := r.config.MLPlatformV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS ML Platform client", err.Error())
		return
	}

	sparkK8SID := data.ID.ValueString()

	err = clusters.Delete(mlPlatformClient, sparkK8SID).ExtractErr()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Unable to delete resource vkcs_mlplatform_spark_k8s", err.Error())
		return
	}

	sparkK8SStateConf := &retry.StateChangeConf{
		Pending:    []string{clusterStatusRunning, clusterStatusDeleting},
		Target:     []string{clusterStatusDeleted},
		Refresh:    clusterStateRefreshFunc(mlPlatformClient, sparkK8SID),
		Timeout:    timeout,
		Delay:      clusterDelay,
		MinTimeout: clusterMinTimeout,
	}

	tflog.Debug(ctx, "Waiting for the Spark K8S cluster to be deleted")

	_, err = sparkK8SStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for the cluster to become ready", err.Error())
		return
	}
}

func (r *SparkK8SResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
