package networking

import (
	"context"
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/dns"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/framework/utils"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/framework/validators"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"golang.org/x/exp/slices"
)

var (
	_ datasource.DataSource              = &PortDataSource{}
	_ datasource.DataSourceWithConfigure = &PortDataSource{}
)

func NewPortDataSource() datasource.DataSource {
	return &PortDataSource{}
}

type PortDataSource struct {
	config clients.Config
}

type PortDataSourceModel struct {
	Region types.String `tfsdk:"region"`
	SDN    types.String `tfsdk:"sdn"`

	AdminStateUp        types.Bool                              `tfsdk:"admin_state_up"`
	AllFixedIPs         types.List                              `tfsdk:"all_fixed_ips"`
	AllSecurityGroupIDs types.Set                               `tfsdk:"all_security_group_ids"`
	AllTags             types.Set                               `tfsdk:"all_tags"`
	AllowedAddressPairs []PortDataSourceAllowedAddressPairModel `tfsdk:"allowed_address_pairs"`
	Description         types.String                            `tfsdk:"description"`
	DeviceID            types.String                            `tfsdk:"device_id"`
	DeviceOwner         types.String                            `tfsdk:"device_owner"`
	DNSAssignment       types.List                              `tfsdk:"dns_assignment"`
	DNSName             types.String                            `tfsdk:"dns_name"`
	ExtraDHCPOption     []PortDataSourceExtraDHCPOptionModel    `tfsdk:"extra_dhcp_option"`
	FixedIP             types.String                            `tfsdk:"fixed_ip"`
	ID                  types.String                            `tfsdk:"id"`
	MACAddress          types.String                            `tfsdk:"mac_address"`
	Name                types.String                            `tfsdk:"name"`
	NetworkID           types.String                            `tfsdk:"network_id"`
	PortID              types.String                            `tfsdk:"port_id"`
	ProjectID           types.String                            `tfsdk:"project_id"`
	SecurityGroupIDs    types.Set                               `tfsdk:"security_group_ids"`
	Status              types.String                            `tfsdk:"status"`
	Tags                types.Set                               `tfsdk:"tags"`
	TenantID            types.String                            `tfsdk:"tenant_id"`
}

type PortDataSourceAllowedAddressPairModel struct {
	IPAddress  types.String `tfsdk:"ip_address"`
	MACAddress types.String `tfsdk:"mac_address"`
}

type PortDataSourceExtraDHCPOptionModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func (d *PortDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_networking_port"
}

func (d *PortDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Optional:    true,
				Description: "The region in which to obtain the Network client. A Network client is needed to retrieve port ids. If omitted, the `region` argument of the provider is used.",
			},

			"sdn": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "SDN to use for this resource. Must be one of following: \"neutron\", \"sprut\". Default value is \"neutron\".",
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("neutron", "sprut"),
				},
			},

			"admin_state_up": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The administrative state of the port.",
			},

			"all_fixed_ips": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "The collection of Fixed IP addresses on the port in the order returned by the Network v2 API.",
			},

			"all_security_group_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "The set of security group IDs applied on the port.",
			},

			"all_tags": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "The set of string tags applied on the port.",
			},

			"allowed_address_pairs": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip_address": schema.StringAttribute{
							Computed:    true,
							Description: "The additional IP address.",
						},

						"mac_address": schema.StringAttribute{
							Computed:    true,
							Description: "The additional MAC address.",
						},
					},
				},
				Computed:    true,
				Description: "An IP/MAC Address pair of additional IP addresses that can be active on this port. The structure is described below.",
			},

			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Human-readable description of the port.",
			},

			"device_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the device the port belongs to.",
			},

			"device_owner": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The device owner of the port.",
			},

			"dns_assignment": schema.ListAttribute{
				ElementType: types.MapType{ElemType: types.StringType},
				Computed:    true,
				Description: "The list of maps representing port DNS assignments.",
			},

			"dns_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The port DNS name to filter.",
			},

			"extra_dhcp_option": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the DHCP option.",
						},

						"value": schema.StringAttribute{
							Computed:    true,
							Description: "Value of the DHCP option.",
						},
					},
				},
				Computed:    true,
				Description: "An extra DHCP option configured on the port. The structure is described below.",
			},

			"fixed_ip": schema.StringAttribute{
				Optional:    true,
				Description: "The port IP address filter.",
				Validators: []validator.String{
					validators.IPAddress(),
				},
			},

			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the found port.",
			},

			"mac_address": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The MAC address of the port.",
			},

			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the port.",
			},

			"network_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the network the port belongs to.",
			},

			"port_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the port.",
			},

			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The project_id of the owner of the port.",
			},

			"security_group_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The list of port security group IDs to filter.",
			},

			"status": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The status of the port.",
			},

			"tags": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The list of port tags to filter.",
			},

			"tenant_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The tenant_id of the owner of the port.",
			},
		},
		Description: "Use this data source to get the ID of an available VKCS port.",
	}
}

func (d *PortDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *PortDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PortDataSourceModel

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
		sdn = networking.DefaultSDN
	}

	client, err := d.config.NetworkingV2Client(region, sdn)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS Networking API client", err.Error())
		return
	}

	var listOpts ports.ListOptsBuilder

	opts := ports.ListOpts{
		Description: data.Description.ValueString(),
		DeviceID:    data.DeviceID.ValueString(),
		DeviceOwner: data.DeviceOwner.ValueString(),
		MACAddress:  data.MACAddress.ValueString(),
		Name:        data.Name.ValueString(),
		NetworkID:   data.NetworkID.ValueString(),
		ID:          data.PortID.ValueString(),
		ProjectID:   data.ProjectID.ValueString(),
		Status:      data.Status.ValueString(),
		TenantID:    data.TenantID.ValueString(),
	}

	if utils.IsKnown(data.AdminStateUp) {
		opts.AdminStateUp = data.AdminStateUp.ValueBoolPointer()
	}

	opts.Tags = expandPortTags(ctx, data.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	listOpts = opts

	if utils.IsKnown(data.DNSName) {
		optsExt := dns.PortListOptsExt{
			ListOptsBuilder: opts,
			DNSName:         data.DNSName.ValueString(),
		}
		listOpts = optsExt
	}

	tflog.Debug(ctx, "Calling Networking API to list ports", map[string]interface{}{"list_opts": fmt.Sprintf("%#v", listOpts)})

	allPages, err := ports.List(client, listOpts).AllPages()
	if err != nil {
		resp.Diagnostics.AddError("Error calling VKCS Networking API", err.Error())
		return
	}

	var allPorts []portExtended
	err = ports.ExtractPortsInto(allPages, &allPorts)
	if err != nil {
		resp.Diagnostics.AddError("Error reading VKCS Networking API response", err.Error())
		return
	}

	tflog.Debug(ctx, "Called Networking API to list ports", map[string]interface{}{"all_ports_len": len(allPorts)})

	if len(allPorts) < 1 {
		resp.Diagnostics.AddError("Your query returned no results",
			"Please change your search criteria and try again")
		return
	}

	tflog.Debug(ctx, "Filtering retrieved ports")

	var filteredPorts []portExtended

	if utils.IsKnown(data.FixedIP) {
		v := data.FixedIP.ValueString()

		for _, p := range allPorts {
			for _, ip := range p.FixedIPs {
				if ip.IPAddress == v {
					filteredPorts = append(filteredPorts, p)
				}
			}
		}
	} else {
		filteredPorts = allPorts
	}

	if utils.IsKnown(data.SecurityGroupIDs) {
		secGroups := expandPortSecurityGroupIDs(ctx, data.SecurityGroupIDs, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		var sgPorts []portExtended
		for _, p := range filteredPorts {
			for _, sg := range p.SecurityGroups {
				if slices.Contains(secGroups, sg) {
					sgPorts = append(sgPorts, p)
				}
			}
		}

		filteredPorts = sgPorts
	}

	tflog.Debug(ctx, "Filtered retrieved ports", map[string]interface{}{"filtered_ports_len": len(filteredPorts)})

	if len(filteredPorts) < 1 {
		resp.Diagnostics.AddError("Your query returned no results",
			"Please change your search criteria and try again")
		return
	}

	if len(filteredPorts) > 1 {
		resp.Diagnostics.AddError("Your query returned more than one result",
			"Please try a more specific search criteria")
		return
	}

	port := filteredPorts[0]
	tflog.Debug(ctx, "Retrieved port", map[string]interface{}{"port": fmt.Sprintf("%#v", port)})

	data.ID = types.StringValue(port.ID)
	data.Region = types.StringValue(region)
	data.SDN = types.StringValue(sdn)

	data.AdminStateUp = types.BoolValue(port.AdminStateUp)
	data.AllFixedIPs = flattenPortAllFixedIPs(ctx, port.FixedIPs, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	data.AllSecurityGroupIDs = flattenAllSecurityGroupIDs(ctx, port.SecurityGroups, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	data.AllTags = flattenPortAllTags(ctx, port.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	data.AllowedAddressPairs = flattenPortAllowedAddressPairs(ctx, port.AllowedAddressPairs, port.MACAddress, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Description = types.StringValue(port.Description)
	data.DeviceID = types.StringValue(port.DeviceID)
	data.DeviceOwner = types.StringValue(port.DeviceOwner)
	data.DNSAssignment = flattenPortDNSAssignment(ctx, port.DNSAssignment, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	data.DNSName = types.StringValue(port.DNSName)
	data.ExtraDHCPOption = flattenPortExtraDHCPOptions(ctx, port.ExtraDHCPOptsExt, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	data.MACAddress = types.StringValue(port.MACAddress)
	data.Name = types.StringValue(port.Name)
	data.NetworkID = types.StringValue(port.NetworkID)
	data.PortID = types.StringValue(port.ID)
	data.ProjectID = types.StringValue(port.ProjectID)
	data.Status = types.StringValue(port.Status)
	data.TenantID = types.StringValue(port.TenantID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenPortAllFixedIPs(ctx context.Context, in []ports.IP, respDiags *diag.Diagnostics) types.List {
	addrs := make([]string, len(in))
	for i, ip := range in {
		addrs[i] = ip.IPAddress
	}

	r, diags := types.ListValueFrom(ctx, types.StringType, addrs)
	respDiags.Append(diags...)
	return r
}

func flattenPortAllTags(ctx context.Context, in []string, respDiags *diag.Diagnostics) types.Set {
	r, diags := types.SetValueFrom(ctx, types.StringType, in)
	respDiags.Append(diags...)
	return r
}

func flattenPortAllowedAddressPairs(ctx context.Context, in []ports.AddressPair, mac string, _ *diag.Diagnostics) []PortDataSourceAllowedAddressPairModel {
	pairs := make([]PortDataSourceAllowedAddressPairModel, len(in))

	for i, pair := range in {
		pairs[i] = PortDataSourceAllowedAddressPairModel{
			IPAddress: types.StringValue(pair.IPAddress),
		}
		// Only set the MAC address if it is different than the
		// port's MAC. This means that a specific MAC was set.
		if pair.MACAddress != mac {
			pairs[i].MACAddress = types.StringValue(pair.MACAddress)
		}
	}

	return pairs
}

func flattenAllSecurityGroupIDs(ctx context.Context, in []string, respDiags *diag.Diagnostics) types.Set {
	r, diags := types.SetValueFrom(ctx, types.StringType, in)
	respDiags.Append(diags...)
	return r
}

func flattenPortDNSAssignment(ctx context.Context, in []map[string]string, respDiags *diag.Diagnostics) types.List {
	r, diags := types.ListValueFrom(ctx, types.MapType{ElemType: types.StringType}, in)
	respDiags.Append(diags...)
	return r
}

func flattenPortExtraDHCPOptions(ctx context.Context, in extradhcpopts.ExtraDHCPOptsExt, _ *diag.Diagnostics) []PortDataSourceExtraDHCPOptionModel {
	r := make([]PortDataSourceExtraDHCPOptionModel, len(in.ExtraDHCPOpts))
	for i, opt := range in.ExtraDHCPOpts {
		r[i] = PortDataSourceExtraDHCPOptionModel{
			Name:  types.StringValue(opt.OptName),
			Value: types.StringValue(opt.OptValue),
		}
	}
	return r
}

func expandPortTags(ctx context.Context, in types.Set, respDiags *diag.Diagnostics) string {
	var tags []string

	respDiags.Append(in.ElementsAs(ctx, &tags, true)...)
	if respDiags.HasError() {
		return ""
	}

	return strings.Join(tags, ",")
}

func expandPortSecurityGroupIDs(ctx context.Context, in types.Set, respDiags *diag.Diagnostics) []string {
	var sgIDs []string
	respDiags.Append(in.ElementsAs(ctx, &sgIDs, true)...)
	return sgIDs
}
