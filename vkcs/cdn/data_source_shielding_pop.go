package cdn

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/cdn/datasource_shielding_pop"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/cdn/v1/shieldingpop"
)

var (
	_ datasource.DataSource              = (*shieldingPopDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*shieldingPopDataSource)(nil)
)

func NewShieldingPopDataSource() datasource.DataSource {
	return &shieldingPopDataSource{}
}

type shieldingPopDataSource struct {
	config clients.Config
}

func (d *shieldingPopDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_shielding_pop"
}

func (d *shieldingPopDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_shielding_pop.ShieldingPopDataSourceSchema(ctx)
}

func (d *shieldingPopDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *shieldingPopDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_shielding_pop.ShieldingPopModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	client, err := d.config.CDNV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating CDN API client", err.Error())
		return
	}

	tflog.Trace(ctx, "Calling CDN API to list origin shielding POPs")

	shieldingPops, err := shieldingpop.List(client, d.config.GetTenantID()).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to list origin shielding POP", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to list origin shielding POPs", map[string]interface{}{"shielding_pops": fmt.Sprintf("%#v", shieldingPops)})

	city := data.City.ValueString()
	country := data.Country.ValueString()
	dc := data.Datacenter.ValueString()

	var filteredShieldingPops []shieldingpop.ShieldingPop
	for _, sp := range shieldingPops {
		if city != "" && sp.City != city {
			continue
		}
		if country != "" && sp.Country != country {
			continue
		}
		if dc != "" && sp.Datacenter != dc {
			continue
		}
		filteredShieldingPops = append(filteredShieldingPops, sp)
	}

	tflog.Trace(ctx, "Filtered available shielding POPs", map[string]interface{}{"filtered_shielding_pops": fmt.Sprintf("%#v", filteredShieldingPops)})

	if len(filteredShieldingPops) < 1 {
		resp.Diagnostics.AddError("Error filtering origin shielding POPs", "Your query returned no results. "+
			"Please change your search criteria and try again.")
		return
	}

	if len(filteredShieldingPops) > 1 {
		resp.Diagnostics.AddError("Error filtering origin shielding POPs", "Your query returned more than one result. "+
			"Please try a more specific search criteria.")
		return
	}

	shieldingPop := shieldingPops[0]

	data.Id = types.Int64Value(int64(shieldingPop.ID))
	data.City = types.StringValue(shieldingPop.City)
	data.Country = types.StringValue(shieldingPop.Country)
	data.Datacenter = types.StringValue(shieldingPop.Datacenter)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
