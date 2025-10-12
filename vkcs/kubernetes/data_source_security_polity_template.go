package kubernetes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v1/securitypolicytemplates"
)

var _ datasource.DataSource = &SecurityPolicyTemplateDataSource{}
var _ datasource.DataSourceWithConfigure = &SecurityPolicyTemplateDataSource{}

func NewSecurityPolicyTemplateDataSource() datasource.DataSource {
	return &SecurityPolicyTemplateDataSource{}
}

type SecurityPolicyTemplateDataSource struct {
	config clients.Config
}

type SecurityPolicyTemplateDataSourceModel struct {
	Region              types.String `tfsdk:"region"`
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	SettingsDescription types.String `tfsdk:"settings_description"`
	Version             types.String `tfsdk:"version"`
	CreatedAt           types.String `tfsdk:"created_at"`
	UpdatedAt           types.String `tfsdk:"updated_at"`
}

func (d *SecurityPolicyTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_kubernetes_security_policy_template"
}

func (d *SecurityPolicyTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The region in which to obtain the service client. If omitted, the `region` argument of the provider is used.",
			},

			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource.",
			},

			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the security policy template.",
			},

			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Description of the security policy template.",
			},

			"settings_description": schema.StringAttribute{
				Computed:    true,
				Description: "Security policy settings description.",
			},

			"version": schema.StringAttribute{
				Computed:    true,
				Description: "Version of the security policy template.",
			},

			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Template creation timestamp",
			},

			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Template update timestamp.",
			},
		},

		Description: "Provides a kubernetes security policy template datasource. This can be used to get information about an VKCS kubernetes security policy template.",
	}
}

func (d *SecurityPolicyTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *SecurityPolicyTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SecurityPolicyTemplateDataSourceModel

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

	tflog.Debug(ctx, "Calling Kubernetes API to get list of security policy templates")

	templates, err := securitypolicytemplates.List(client).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling VKCS Kubernetes API", err.Error())
		return
	}

	tflog.Debug(ctx, "Called Kubernetes API to get list of cluster templates", map[string]interface{}{"templates": fmt.Sprintf("%#v", templates)})

	var foundTemplate *securitypolicytemplates.SecurityPolicyTemplate
	for _, template := range templates {
		if template.Name == data.Name.ValueString() {
			foundTemplate = &template
			break
		}
	}

	if foundTemplate == nil {
		resp.Diagnostics.AddError("Error retrieving vkcs_kubernetes_security_policy_template", fmt.Sprintf("Template %s not found", data.Name.ValueString()))
		return
	}

	data.Region = types.StringValue(region)
	data.ID = types.StringValue(foundTemplate.UUID)
	data.Name = types.StringValue(foundTemplate.Name)
	data.Description = types.StringValue(foundTemplate.Description)
	data.SettingsDescription = types.StringValue(foundTemplate.SettingsDescription)
	data.Version = types.StringValue(foundTemplate.Version)
	data.CreatedAt = types.StringValue(foundTemplate.CreatedAt)
	data.UpdatedAt = types.StringValue(foundTemplate.UpdatedAt)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
