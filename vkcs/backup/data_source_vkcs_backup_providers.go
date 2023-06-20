package backup

import (
	"context"

	"github.com/gophercloud/utils/terraform/hashcode"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/providers"
)

// Ensure the implementation satisfies the desired interfaces.
var _ datasource.DataSource = &ProvidersDataSource{}

func NewProvidersDataSource() datasource.DataSource {
	return &ProvidersDataSource{}
}

type ProvidersDataSource struct {
	config clients.Config
}

func (d *ProvidersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_backup_providers"
}

type ProvidersDataSourceModel struct {
	Providers []ProvidersDataSourceProvidersModel `tfsdk:"providers"`
	Region    types.String                        `tfsdk:"region"`
	ID        types.String                        `tfsdk:"id"`
}

type ProvidersDataSourceProvidersModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *ProvidersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"providers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the backup provider",
						},
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Name of the backup provider",
						},
					},
				},
				Description: "List of available backup providers",
			},

			"region": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The `region` to fetch availability zones from, defaults to the provider's `region`.",
			},

			"id": schema.StringAttribute{
				Computed: true,
			},
		},
		Description: "Use this data source to get backup providers info",
	}
}

func (d *ProvidersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *ProvidersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProvidersDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}

	backupClient, err := d.config.BackupV1Client(region, d.config.GetTenantID())
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS backup client", err.Error())
		return
	}

	providersInfo, err := providers.List(backupClient).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vkcs_backup_providers", err.Error())
		return
	}

	flattenedProviders := flattenProviders(providersInfo)

	var names []string
	for _, p := range flattenedProviders {
		names = append(names, p.Name.ValueString())
	}

	data.Providers = flattenedProviders
	data.Region = types.StringValue(region)
	data.ID = types.StringValue(hashcode.Strings(names))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenProviders(providersInfo []*providers.Provider) (r []ProvidersDataSourceProvidersModel) {
	for _, p := range providersInfo {
		r = append(r, ProvidersDataSourceProvidersModel{
			ID:   types.StringValue(p.ID),
			Name: types.StringValue(providerNameMapping[p.Name]),
		})
	}
	return
}
