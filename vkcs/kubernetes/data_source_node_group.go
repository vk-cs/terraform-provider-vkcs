package kubernetes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/framework/utils"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v1/nodegroups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

var (
	_ datasource.DataSource              = &NodeGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &NodeGroupDataSource{}
)

func NewNodeGroupDataSource() datasource.DataSource {
	return &NodeGroupDataSource{}
}

type NodeGroupDataSource struct {
	config clients.Config
}

type NodeGroupDataSourceModel struct {
	ID     types.String `tfsdk:"id"`
	Region types.String `tfsdk:"region"`

	AutoscalingEnabled types.Bool           `tfsdk:"autoscaling_enabled"`
	AvailabilityZones  types.List           `tfsdk:"availability_zones"`
	ClusterID          types.String         `tfsdk:"cluster_id"`
	FlavorID           types.String         `tfsdk:"flavor_id"`
	MaxNodeUnavailable types.Int64          `tfsdk:"max_node_unavailable"`
	MaxNodes           types.Int64          `tfsdk:"max_nodes"`
	MinNodes           types.Int64          `tfsdk:"min_nodes"`
	Name               types.String         `tfsdk:"name"`
	NodeCount          types.Int64          `tfsdk:"node_count"`
	Nodes              []NodeGroupNodeModel `tfsdk:"nodes"`
	State              types.String         `tfsdk:"state"`
	UUID               types.String         `tfsdk:"uuid"`
	VolumeSize         types.Int64          `tfsdk:"volume_size"`
	VolumeType         types.String         `tfsdk:"volume_type"`
}

type NodeGroupNodeModel struct {
	CreatedAt   types.String `tfsdk:"created_at"`
	Name        types.String `tfsdk:"name"`
	NodeGroupID types.String `tfsdk:"node_group_id"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	UUID        types.String `tfsdk:"uuid"`
}

func (d *NodeGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_kubernetes_node_group"
}

func (d *NodeGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The UUID of the cluster's node group.",
			},

			"region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The region to obtain the service client. If omitted, the `region` argument of the provider is used.",
			},

			"autoscaling_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Determines whether the autoscaling is enabled.",
			},

			"availability_zones": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "The list of availability zones of the node group.",
			},

			"cluster_id": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of cluster that node group belongs.",
			},

			"flavor_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The id of the flavor.",
			},

			"max_node_unavailable": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Specified as a percentage. The maximum number of nodes that can fail during an upgrade.",
			},

			"max_nodes": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The maximum amount of nodes in the node group.",
			},

			"min_nodes": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The minimum amount of nodes in the node group.",
			},

			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the node group.",
			},

			"node_count": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The count of nodes in the node group.",
			},

			"nodes": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "Time when a node was created.",
						},

						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of a node.",
						},

						"node_group_id": schema.StringAttribute{
							Computed:    true,
							Description: "The node group id.",
						},

						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "Time when a node was updated.",
						},

						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "UUID of a node.",
						},
					},
				},
				Description: "The list of node group's node objects.",
			},

			"state": schema.StringAttribute{
				Computed:    true,
				Description: "Determines current state of node group (RUNNING, SHUTOFF, ERROR).",
			},

			"uuid": schema.StringAttribute{
				Optional:           true,
				Computed:           true,
				Description:        "The UUID of the cluster's node group.",
				DeprecationMessage: "This argument is deprecated, please, use the `id` attribute instead.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("id"),
					),
				},
			},

			"volume_size": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The amount of memory in the volume in GB",
			},

			"volume_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The type of the volume.",
			},
		},
		Description: "Use this data source to get information on VKCS Kubernetes cluster's node group.",
	}
}

func (d *NodeGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *NodeGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NodeGroupDataSourceModel
	var diags diag.Diagnostics

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}

	client, err := d.config.ContainerInfraV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS Kubernetes API client", err.Error())
		return
	}

	id := utils.GetFirstNotEmptyValue(data.ID, data.UUID)
	ctx = tflog.SetField(ctx, "id", id)

	tflog.Debug(ctx, "Calling Kubernetes API to get the node group")

	nGroup, err := nodegroups.Get(client, id).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling VKCS Kubernetes API", err.Error())
		return
	}

	tflog.Debug(ctx, "Called Kubernetes API to get the node group", map[string]any{"node_group": fmt.Sprintf("%#v", nGroup)})

	data.ID = types.StringValue(id)
	data.UUID = data.ID
	data.Region = types.StringValue(region)
	data.AutoscalingEnabled = types.BoolValue(nGroup.Autoscaling)
	data.AvailabilityZones, diags = types.ListValueFrom(ctx, types.StringType, nGroup.AvailabilityZones)
	resp.Diagnostics.Append(diags...)
	data.ClusterID = types.StringValue(nGroup.ClusterID)
	data.FlavorID = types.StringValue(nGroup.FlavorID)
	data.MaxNodeUnavailable = types.Int64Value(int64(nGroup.MaxNodeUnavailable))
	data.MaxNodes = types.Int64Value(int64(nGroup.MaxNodes))
	data.MinNodes = types.Int64Value(int64(nGroup.MinNodes))
	data.Name = types.StringValue(nGroup.Name)
	data.NodeCount = types.Int64Value(int64(nGroup.NodeCount))
	data.Nodes = flattenNodes(nGroup.Nodes)
	data.State = types.StringValue(nGroup.State)
	data.UUID = types.StringValue(nGroup.UUID)
	data.VolumeSize = types.Int64Value(int64(nGroup.VolumeSize))
	data.VolumeType = types.StringValue(nGroup.VolumeType)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenNodes(nodes []*nodegroups.Node) (r []NodeGroupNodeModel) {
	for _, n := range nodes {
		if n == nil {
			continue
		}
		r = append(r, NodeGroupNodeModel{
			CreatedAt:   types.StringValue(util.GetTimestamp(n.CreatedAt)),
			Name:        types.StringValue(n.Name),
			NodeGroupID: types.StringValue(n.NodeGroupID),
			UpdatedAt:   types.StringValue(util.GetTimestamp(n.UpdatedAt)),
			UUID:        types.StringValue(n.UUID),
		})
	}
	return
}
