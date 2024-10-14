package cdn

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/cdn/datasource_origin_group"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/cdn/v1/origingroups"
)

var (
	_ datasource.DataSource              = (*originGroupDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*originGroupDataSource)(nil)
)

func NewOriginGroupDataSource() datasource.DataSource {
	return &originGroupDataSource{}
}

type originGroupDataSource struct {
	config clients.Config
}

func (d *originGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdn_origin_group"
}

func (d *originGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_origin_group.OriginGroupDataSourceSchema(ctx)
}

func (d *originGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *originGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_origin_group.OriginGroupModel

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

	tflog.Trace(ctx, "Calling CDN API to list origin groups")

	originGroups, err := origingroups.List(client, d.config.GetTenantID()).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling CDN API to list origin groups", err.Error())
		return
	}

	tflog.Trace(ctx, "Called CDN API to list origin groups", map[string]interface{}{"origin_groups": fmt.Sprintf("%#v", originGroups)})

	name := data.Name.ValueString()
	i := slices.IndexFunc(originGroups, func(og origingroups.OriginGroup) bool {
		return og.Name == name
	})
	if i == -1 {
		resp.Diagnostics.AddError("Error finding an origin group", fmt.Sprintf("Origin group with name %q not found", name))
		return
	}

	originGroup := originGroups[i]

	data.Id = types.Int64Value(int64(originGroup.ID))
	data.Name = types.StringValue(originGroup.Name)

	var diags diag.Diagnostics
	data.Origins, diags = datasource_origin_group.FlattenOrigins(ctx, originGroup.Origins)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.UseNext = types.BoolValue(originGroup.UseNext)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
