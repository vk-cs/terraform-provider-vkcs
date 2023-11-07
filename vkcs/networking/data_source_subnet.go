package networking

import (
	"context"
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/framework/utils"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	isubnets "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking/v2/subnets"
)

var (
	_ datasource.DataSource              = &SubnetDataSource{}
	_ datasource.DataSourceWithConfigure = &SubnetDataSource{}
)

func NewSubnetDataSource() datasource.DataSource {
	return &SubnetDataSource{}
}

type SubnetDataSource struct {
	config clients.Config
}

type SubnetDataSourceModel struct {
	SDN    types.String `tfsdk:"sdn"`
	Region types.String `tfsdk:"region"`

	AllTags         types.Set                             `tfsdk:"all_tags"`
	AllocationPools []SubnetDataSourceAllocationPoolModel `tfsdk:"allocation_pools"`
	CIDR            types.String                          `tfsdk:"cidr"`
	Description     types.String                          `tfsdk:"description"`
	DHCPEnabled     types.Bool                            `tfsdk:"dhcp_enabled"`
	DNSNameservers  types.Set                             `tfsdk:"dns_nameservers"`
	EnableDHCP      types.Bool                            `tfsdk:"enable_dhcp"`
	GatewayIP       types.String                          `tfsdk:"gateway_ip"`
	HostRoutes      []SubnetDataSourceHostRouteModel      `tfsdk:"host_routes"`
	ID              types.String                          `tfsdk:"id"`
	Name            types.String                          `tfsdk:"name"`
	NetworkID       types.String                          `tfsdk:"network_id"`
	SubnetID        types.String                          `tfsdk:"subnet_id"`
	SubnetPoolID    types.String                          `tfsdk:"subnetpool_id"`
	Tags            types.Set                             `tfsdk:"tags"`
	TenantID        types.String                          `tfsdk:"tenant_id"`
}

type SubnetDataSourceAllocationPoolModel struct {
	End   types.String `tfsdk:"end"`
	Start types.String `tfsdk:"start"`
}

type SubnetDataSourceHostRouteModel struct {
	DestinationCIDR types.String `tfsdk:"destination_cidr"`
	NextHop         types.String `tfsdk:"next_hop"`
}

func (d *SubnetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_networking_subnet"
}

func (d *SubnetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The region in which to obtain the Network client. A Network client is needed to retrieve subnet ids. If omitted, the `region` argument of the provider is used.",
			},

			"sdn": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "SDN to use for this resource. Must be one of following: \"neutron\", \"sprut\". Default value is project's default SDN.",
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("neutron", "sprut"),
				},
			},

			"all_tags": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "A set of string tags applied on the subnet.",
			},

			"allocation_pools": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"end": schema.StringAttribute{
							Computed:    true,
							Description: "The ending address.",
						},

						"start": schema.StringAttribute{
							Computed:    true,
							Description: "The starting address.",
						},
					},
				},
				Computed:    true,
				Description: "Allocation pools of the subnet.",
			},

			"cidr": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The CIDR of the subnet.",
			},

			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Human-readable description of the subnet.",
			},

			"dhcp_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "If the subnet has DHCP enabled.",
			},

			"dns_nameservers": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "DNS Nameservers of the subnet.",
			},

			"enable_dhcp": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the subnet has DHCP enabled or not.",
			},

			"gateway_ip": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The IP of the subnet's gateway.",
			},

			"host_routes": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"destination_cidr": schema.StringAttribute{
							Computed: true,
						},

						"next_hop": schema.StringAttribute{
							Computed: true,
						},
					},
				},
				Computed:    true,
				Description: "Host Routes of the subnet.",
			},

			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the found subnet.",
			},

			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the subnet.",
			},

			"network_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the network the subnet belongs to.",
			},

			"subnet_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the subnet.",
			},

			"subnetpool_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the subnetpool associated with the subnet.",
			},

			"tags": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The list of subnet tags to filter.",
			},

			"tenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The owner of the subnet.",
			},
		},
		Description: "Use this data source to get the ID of an available VKCS subnet.",
	}
}

func (d *SubnetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *SubnetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SubnetDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}

	sdn := data.SDN.ValueString()
	if sdn == "" {
		sdn = networking.SearchInAllSDNs
	}

	client, err := d.config.NetworkingV2Client(region, sdn)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS Networking API client", err.Error())
		return
	}

	listOpts := subnets.ListOpts{
		Name:         data.Name.ValueString(),
		Description:  data.Description.ValueString(),
		NetworkID:    data.NetworkID.ValueString(),
		TenantID:     data.TenantID.ValueString(),
		GatewayIP:    data.GatewayIP.ValueString(),
		CIDR:         data.CIDR.ValueString(),
		ID:           data.SubnetID.ValueString(),
		SubnetPoolID: data.SubnetPoolID.ValueString(),
	}

	if utils.IsKnown(data.DHCPEnabled) {
		listOpts.EnableDHCP = data.DHCPEnabled.ValueBoolPointer()
	}

	listOpts.Tags = expandSubnetDataSourceTags(ctx, data.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Calling Networking API to list subnets", map[string]interface{}{"list_opts": fmt.Sprintf("%#v", listOpts)})

	allPages, err := subnets.List(client, &listOpts).AllPages()
	if err != nil {
		resp.Diagnostics.AddError("Error calling VKCS Networking API", err.Error())
		return
	}

	var allSubnets []subnetExtended
	err = isubnets.ExtractSubnetsInto(allPages, &allSubnets)
	if err != nil {
		resp.Diagnostics.AddError("Error reading VKCS Networking API response", err.Error())
		return
	}

	tflog.Debug(ctx, "Called Networking API to list subnets", map[string]interface{}{"all_subnets_len": len(allSubnets)})

	if len(allSubnets) < 1 {
		resp.Diagnostics.AddError("Your query returned no results",
			"Please change your search criteria and try again")
		return
	}

	if len(allSubnets) > 1 {
		resp.Diagnostics.AddError("Your query returned more than one result",
			"Please try a more specific search criteria")
		return
	}

	subnet := allSubnets[0]
	tflog.Debug(ctx, "Retrieved the subnet", map[string]interface{}{"subnet": fmt.Sprintf("%#v", subnet)})

	data.ID = types.StringValue(subnet.ID)
	data.Region = types.StringValue(region)
	data.SDN = types.StringValue(subnet.SDN)

	data.AllTags = flattenSubnetDataSourceAllTags(ctx, subnet.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	data.AllocationPools = flattenSubnetDataSourceAllocationPools(ctx, subnet.AllocationPools, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	data.CIDR = types.StringValue(subnet.CIDR)
	data.Description = types.StringValue(subnet.Description)
	data.DNSNameservers = flattenSubnetDataSourceDNSNameservers(ctx, subnet.DNSNameservers, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	data.EnableDHCP = types.BoolValue(subnet.EnableDHCP)
	data.GatewayIP = types.StringValue(subnet.GatewayIP)
	data.HostRoutes = flattenSubnetDataSourceHostRoutes(ctx, subnet.HostRoutes, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Name = types.StringValue(subnet.Name)
	data.NetworkID = types.StringValue(subnet.NetworkID)
	data.SubnetPoolID = types.StringValue(subnet.SubnetPoolID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenSubnetDataSourceAllTags(ctx context.Context, in []string, respDiags *diag.Diagnostics) types.Set {
	r, diags := types.SetValueFrom(ctx, types.StringType, in)
	respDiags.Append(diags...)
	return r
}

func flattenSubnetDataSourceAllocationPools(_ context.Context, in []subnets.AllocationPool, _ *diag.Diagnostics) []SubnetDataSourceAllocationPoolModel {
	r := make([]SubnetDataSourceAllocationPoolModel, len(in))
	for i, ap := range in {
		r[i] = SubnetDataSourceAllocationPoolModel{
			End:   types.StringValue(ap.End),
			Start: types.StringValue(ap.Start),
		}
	}
	return r
}

func flattenSubnetDataSourceDNSNameservers(ctx context.Context, in []string, respDiags *diag.Diagnostics) types.Set {
	r, diags := types.SetValueFrom(ctx, types.StringType, in)
	respDiags.Append(diags...)
	return r
}

func flattenSubnetDataSourceHostRoutes(_ context.Context, in []subnets.HostRoute, _ *diag.Diagnostics) []SubnetDataSourceHostRouteModel {
	r := make([]SubnetDataSourceHostRouteModel, len(in))
	for i, hr := range in {
		r[i] = SubnetDataSourceHostRouteModel{
			DestinationCIDR: types.StringValue(hr.DestinationCIDR),
			NextHop:         types.StringValue(hr.NextHop),
		}
	}
	return r
}

func expandSubnetDataSourceTags(ctx context.Context, in types.Set, respDiags *diag.Diagnostics) string {
	var tags []string
	respDiags.Append(in.ElementsAs(ctx, &tags, true)...)
	return strings.Join(tags, "")
}
