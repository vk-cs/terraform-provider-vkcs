package dc

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dc/v2/apioptions"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
)

// Ensure the implementation satisfies the desired interfaces.
var _ datasource.DataSource = &APIOptionsDataSource{}

func NewAPIOptionsDataSource() datasource.DataSource {
	return &APIOptionsDataSource{}
}

type APIOptionsDataSource struct {
	config clients.Config
}

func (d *APIOptionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_dc_api_options"
}

type APIOptionDataSourceModel struct {
	ID                types.String   `tfsdk:"id"`
	AvailabilityZones []types.String `tfsdk:"availability_zones"`
	Flavors           []types.String `tfsdk:"flavors"`
	Region            types.String   `tfsdk:"region"`
}

func (d *APIOptionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the data source",
			},

			"availability_zones": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of avalability zone options",
			},

			"flavors": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of flavor options for vkcs_dc_router resource",
			},

			"region": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The `region` to fetch availability zones from, defaults to the provider's `region`.",
			},
		},
		Description: "Use this data source to get direct connect api options._note_ This data source requires Sprut SDN to be enabled in your project.",
	}
}

func (d *APIOptionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *APIOptionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data APIOptionDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}

	networkingClient, err := d.config.NetworkingV2Client(region, networking.SprutSDN)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS networking client", err.Error())
		return
	}

	apiOptionsResp, err := apioptions.Get(networkingClient).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vkcs_dc_api_options", err.Error())
		return
	}

	data.ID = types.StringValue(uuid.New().String())

	availabilityZones := make([]types.String, len(apiOptionsResp.AvailabilityZones))
	for i, az := range apiOptionsResp.AvailabilityZones {
		availabilityZones[i] = types.StringValue(az)
	}
	data.AvailabilityZones = availabilityZones

	flavors := make([]types.String, len(apiOptionsResp.Flavors))
	for i, flavor := range apiOptionsResp.Flavors {
		flavors[i] = types.StringValue(flavor)
	}
	data.Flavors = flavors

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
