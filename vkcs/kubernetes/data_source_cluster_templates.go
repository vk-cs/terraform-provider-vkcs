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
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v1/clustertemplates"
)

var (
	_ datasource.DataSource              = &ClusterTemplatesDataSource{}
	_ datasource.DataSourceWithConfigure = &ClusterTemplatesDataSource{}
)

func NewClusterTemplatesDataSource() datasource.DataSource {
	return &ClusterTemplatesDataSource{}
}

type ClusterTemplatesDataSource struct {
	config clients.Config
}

type ClusterTemplatesDataSourceModel struct {
	Region types.String `tfsdk:"region"`

	ClusterTemplates []ClusterTemplateModel `tfsdk:"cluster_templates"`
	ID               types.String           `tfsdk:"id"`
}

type ClusterTemplateModel struct {
	ClusterTemplateUUID types.String `tfsdk:"cluster_template_uuid"`
	Name                types.String `tfsdk:"name"`
	Version             types.String `tfsdk:"version"`
}

func (d *ClusterTemplatesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_kubernetes_clustertemplates"
}

func (d *ClusterTemplatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The region to obtain the service client. If omitted, the `region` argument of the provider is used.",
			},

			"cluster_templates": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cluster_template_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "UUID of a cluster template.",
						},

						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of a cluster template.",
						},

						"version": schema.StringAttribute{
							Computed:    true,
							Description: "Version of Kubernetes.",
						},
					},
				},
				Description: "Available kubernetes cluster templates.",
			},

			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Random identifier of the data source.",
			},
		},
		Description: "Use this data source to get a list of available VKCS Kubernetes Cluster Templates. To get details about each cluster template the data source can be combined with the `vkcs_kubernetes_clustertemplate` data source.",
	}
}

func (d *ClusterTemplatesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *ClusterTemplatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ClusterTemplatesDataSourceModel

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

	tflog.Debug(ctx, "Calling Kubernetes API to get list of cluster templates")

	templates, err := clustertemplates.List(client).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling VKCS Kubernetes API", err.Error())
		return
	}

	tflog.Debug(ctx, "Called Kubernetes API to get list of cluster templates", map[string]interface{}{"templates": fmt.Sprintf("%#v", templates)})

	data.Region = types.StringValue(region)
	data.ClusterTemplates = flattenClusterTemplates(templates)
	data.ID = types.StringValue(strconv.FormatInt(time.Now().Unix(), 10))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenClusterTemplates(templates []clustertemplates.ClusterTemplate) (r []ClusterTemplateModel) {
	for _, t := range templates {
		r = append(r, ClusterTemplateModel{
			ClusterTemplateUUID: types.StringValue(t.UUID),
			Name:                types.StringValue(t.Name),
			Version:             types.StringValue(t.Version),
		})
	}
	return
}
