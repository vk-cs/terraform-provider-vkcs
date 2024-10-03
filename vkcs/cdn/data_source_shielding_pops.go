package cdn

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/cdn/datasource_shielding_pops"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/cdn/v1/shieldingpop"
)

var (
	_ datasource.DataSource              = (*shieldingPopsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*shieldingPopsDataSource)(nil)
)

func NewShieldingPopsDataSource() datasource.DataSource {
	return &shieldingPopsDataSource{}
}

type shieldingPopsDataSource struct {
	config clients.Config
}

func (d *shieldingPopsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_shielding_pops"
}

func (d *shieldingPopsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_shielding_pops.ShieldingPopsDataSourceSchema(ctx)
}

func (d *shieldingPopsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *shieldingPopsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_shielding_pops.ShieldingPopsModel

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

	data.Id = types.StringValue(strconv.FormatInt(time.Now().Unix(), 10))

	var diags diag.Diagnostics
	data.ShieldingPops, diags = datasource_shielding_pops.FlattenShieldingePops(ctx, shieldingPops)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
