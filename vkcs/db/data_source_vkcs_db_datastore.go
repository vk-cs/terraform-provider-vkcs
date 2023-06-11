package db

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/datastores"
)

var (
	_ datasource.DataSource              = &DatastoreDataSource{}
	_ datasource.DataSourceWithConfigure = &DatastoreDataSource{}
)

func NewDatastoreDataSource() datasource.DataSource {
	return &DatastoreDataSource{}
}

type DatastoreDataSource struct {
	config clients.Config
}

type DatastoreDataSourceModel struct {
	Region types.String `tfsdk:"region"`

	ClusterVolumeTypes types.List              `tfsdk:"cluster_volume_types"`
	ID                 types.String            `tfsdk:"id"`
	MinimumCPU         types.Int64             `tfsdk:"minimum_cpu"`
	MinimumRAM         types.Int64             `tfsdk:"minimum_ram"`
	Name               types.String            `tfsdk:"name"`
	Versions           []DatastoreVersionModel `tfsdk:"versions"`
	VolumeTypes        types.List              `tfsdk:"volume_types"`
}

type DatastoreVersionModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *DatastoreDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_db_datastore"
}

func (d *DatastoreDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The region to obtain the service client. If omitted, the `region` argument of the provider is used.",
			},

			"cluster_volume_types": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Supported volume types for the datastore when used in a cluster.",
			},

			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The id of the datastore.",
			},

			"minimum_cpu": schema.Int64Attribute{
				Computed:    true,
				Description: "Minimum CPU required for instance of the datastore.",
			},

			"minimum_ram": schema.Int64Attribute{
				Computed:    true,
				Description: "Minimum RAM required for instance of the datastore.",
			},

			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the datastore.",
			},

			"versions": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of a version of the datastore.",
						},

						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of a version of the datastore.",
						},
					},
				},
				Description: "Versions of the datastore.",
			},

			"volume_types": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Supported volume types for the datastore.",
			},
		},
		Description: "Use this data source to get information on a VKCS db datastore.",
	}
}

func (d *DatastoreDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *DatastoreDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DatastoreDataSourceModel
	var diags diag.Diagnostics

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

	id, name := data.ID.ValueString(), data.Name.ValueString()
	ctx = tflog.SetField(ctx, "id", id)
	ctx = tflog.SetField(ctx, "name", name)

	tflog.Debug(ctx, "Filtering retrieved datastores")
	filteredDatastores := filterDatastores(allDatastores, id, name)
	tflog.Debug(ctx, "Filtered retrieved datastores", map[string]interface{}{"filtered_datastores": fmt.Sprintf("%#v", filteredDatastores)})

	if len(filteredDatastores) < 1 {
		resp.Diagnostics.AddError("Your query returned no results", "Please change your search criteria and try again.")
		return
	}

	if len(filteredDatastores) > 1 {
		resp.Diagnostics.AddError("Your query returned more than one result", "Please try a more specific search criteria")
		return
	}

	id = filteredDatastores[0].ID
	ctx = tflog.SetField(ctx, "id", id)

	tflog.Debug(ctx, "Calling Databases API to get the datastore by id")

	dStore, err := datastores.Get(client, id).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling VKCS Databases API", err.Error())
		return
	}

	tflog.Debug(ctx, "Called Databases API to get the datastore by its id", map[string]interface{}{"datastore": fmt.Sprintf("%#v", dStore)})

	versions := flattenDatastoreVersions(dStore.Versions)
	sort.SliceStable(versions, func(i, j int) bool {
		return versions[i].Name.ValueString() > versions[j].Name.ValueString()
	})

	data.Region = types.StringValue(region)
	data.ClusterVolumeTypes, diags = types.ListValueFrom(ctx, types.StringType, dStore.ClusterVolumeTypes)
	resp.Diagnostics.Append(diags...)
	data.ID = types.StringValue(dStore.ID)
	data.MinimumCPU = types.Int64Value(int64(dStore.MinimumCPU))
	data.MinimumRAM = types.Int64Value(int64(dStore.MinimumRAM))
	data.Name = types.StringValue(dStore.Name)
	data.Versions = versions
	data.VolumeTypes, diags = types.ListValueFrom(ctx, types.StringType, dStore.VolumeTypes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func filterDatastores(dsSlice []datastores.Datastore, id, name string) []datastores.Datastore {
	var res []datastores.Datastore
	for _, ds := range dsSlice {
		if (name == "" || ds.Name == name) && (id == "" || ds.ID == id) {
			res = append(res, ds)
		}
	}
	return res
}

func flattenDatastoreVersions(versions []datastores.Version) (r []DatastoreVersionModel) {
	for _, v := range versions {
		r = append(r, DatastoreVersionModel{
			ID:   types.StringValue(v.ID),
			Name: types.StringValue(v.Name),
		})
	}
	return
}
