package db

import (
	"context"
	"sort"

	"github.com/gophercloud/utils/terraform/hashcode"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/datastores"
)

var (
	_ datasource.DataSource              = &DatastoresDataSource{}
	_ datasource.DataSourceWithConfigure = &DatastoresDataSource{}
)

func NewDatastoresDataSource() datasource.DataSource {
	return &DatastoresDataSource{}
}

type DatastoresDataSource struct {
	config clients.Config
}

type DatastoresDataSourceModel struct {
	ID     types.String `tfsdk:"id"`
	Region types.String `tfsdk:"region"`

	Datastores []DatastoreModel `tfsdk:"datastores"`
}

type DatastoreModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *DatastoresDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_db_datastores"
}

func (d *DatastoresDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "ID of the resource.",
			},

			"region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The region to obtain the service client. If omitted, the `region` argument of the provider is used.",
			},

			"datastores": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of a datastore.",
						},

						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of a datastore.",
						},
					},
				},
				Description: "List of datastores within VKCS.",
			},
		},
		Description: "Use this data source to get a list of datastores from VKCS.",
	}
}

func (d *DatastoresDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *DatastoresDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DatastoresDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}

	client, err := d.config.DatabaseV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS Databases API client", err.Error())
		return
	}

	tflog.Debug(ctx, "Calling Databases API to list all datastores")

	allPages, err := datastores.List(client).AllPages()
	if err != nil {
		resp.Diagnostics.AddError("Error calling VKCS Databases API", err.Error())
		return
	}

	allDatastores, err := datastores.ExtractDatastores(allPages)
	if err != nil {
		resp.Diagnostics.AddError("Error reading VKCS Databases API response", err.Error())
		return
	}

	tflog.Debug(ctx, "Called Databases API to list all datastores")

	flattenedDatastores := flattenDatastores(allDatastores)
	sort.SliceStable(flattenedDatastores, func(i, j int) bool {
		return flattenedDatastores[i].Name.ValueString() < flattenedDatastores[j].Name.ValueString()
	})

	var names []string
	for _, d := range flattenedDatastores {
		names = append(names, d.Name.ValueString())
	}

	data.ID = types.StringValue(hashcode.Strings(names))
	data.Region = types.StringValue(region)
	data.Datastores = flattenedDatastores

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenDatastores(datastores []datastores.Datastore) (r []DatastoreModel) {
	for _, d := range datastores {
		r = append(r, DatastoreModel{
			ID:   types.StringValue(d.ID),
			Name: types.StringValue(d.Name),
		})
	}
	return
}
