package kubernetes

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v1/securitypolicytemplates"
)

var _ datasource.DataSource = &SecurityPolicyTemplatesDataSource{}
var _ datasource.DataSourceWithConfigure = &SecurityPolicyTemplatesDataSource{}

func NewSecurityPolicyTemplatesDataSource() datasource.DataSource {
	return &SecurityPolicyTemplatesDataSource{}
}

type SecurityPolicyTemplatesDataSource struct {
	config clients.Config
}

type SecurityPolicyTemplatesDataSourceModel struct {
	Region types.String `tfsdk:"region"`

	SecurityPolicyTemplates []SecurityPolicyTemplateModel `tfsdk:"security_policy_templates"`
	ID                      types.String                  `tfsdk:"id"`
}

type SecurityPolicyTemplateModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	SettingsDescription types.String `tfsdk:"settings_description"`
	Version             types.String `tfsdk:"version"`
	CreatedAt           types.String `tfsdk:"created_at"`
	UpdatedAt           types.String `tfsdk:"updated_at"`
}

func (d *SecurityPolicyTemplatesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_kubernetes_security_policy_templates"
}

func (d *SecurityPolicyTemplatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The region in which to obtain the service client. If omitted, the `region` argument of the provider is used.",
			},

			"security_policy_templates": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the template.",
						},

						"name": schema.StringAttribute{
							Computed:    true,
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
				},
				Description: "Available kubernetes security policy templates.",
			},

			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Random identifier of the data source.",
			},
		},

		Description: "Provides a kubernetes security policy templates datasource. This can be used to get information about all available VKCS kubernetes security policy templates.",
	}
}

func (d *SecurityPolicyTemplatesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *SecurityPolicyTemplatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SecurityPolicyTemplatesDataSourceModel

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

	data.Region = types.StringValue(region)
	data.SecurityPolicyTemplates = flattenSecurityPolicyTemplates(templates)
	data.ID = types.StringValue(strconv.FormatInt(time.Now().Unix(), 10))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenSecurityPolicyTemplates(templates []securitypolicytemplates.SecurityPolicyTemplate) (r []SecurityPolicyTemplateModel) {
	for _, t := range templates {
		r = append(r, SecurityPolicyTemplateModel{
			ID:                  types.StringValue(t.UUID),
			Name:                types.StringValue(t.Name),
			Description:         types.StringValue(t.Description),
			SettingsDescription: types.StringValue(t.SettingsDescription),
			Version:             types.StringValue(t.Version),
			CreatedAt:           types.StringValue(t.CreatedAt),
			UpdatedAt:           types.StringValue(t.UpdatedAt),
		})
	}
	return
}
