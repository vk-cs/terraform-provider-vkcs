package dataplatform

import (
	"context"
	"fmt"

	"github.com/vk-cs/terraform-provider-vkcs/vkcs/dataplatform/datasource_template"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/templates"
)

var (
	_ datasource.DataSource              = (*templateDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*templateDataSource)(nil)
)

func NewTemplateDataSource() datasource.DataSource {
	return &templateDataSource{}
}

type templateDataSource struct {
	config clients.Config
}

func (d *templateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dataplatform_template"
}

func (d *templateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_template.TemplateDataSourceSchema(ctx)
}

func (d *templateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *templateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_template.TemplateModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := d.config.GetRegion()

	client, err := d.config.DataPlatformClient(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS dataplatform client", err.Error())
		return
	}

	tflog.Trace(ctx, "Calling Data Platform API to list templates")

	templatesResp, err := templates.Get(client).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling Data Platform API to list products", err.Error())
		return
	}

	tflog.Trace(ctx, "Called Data Platform API to list templates", map[string]interface{}{"templates": fmt.Sprintf("%#v", templatesResp.ClusterTemplates)})

	var templates []templates.ClusterTemplate
	productName := data.ProductName.ValueString()
	productVersion := data.ProductVersion.ValueString()

	for _, template := range templatesResp.ClusterTemplates {
		if template.ProductName == productName {
			if productVersion == "" || template.ProductVersion == productVersion {
				templates = append(templates, template)
			}
		}
	}

	if len(templates) < 1 {
		resp.Diagnostics.AddError("Your query returned no results", "Please change your search criteria and try again.")
		return
	}

	if len(templates) > 1 {
		resp.Diagnostics.AddError("Your query returned more than one result", "Please try a more specific search criteria")
		return
	}

	resp.State.SetAttribute(ctx, path.Root("id"), templates[0].ID)

	diags := data.UpdateFromTemplate(ctx, &templates[0])
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
