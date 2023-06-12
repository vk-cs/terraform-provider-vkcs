package db

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/datastores"
)

var (
	_ datasource.DataSource              = &DatastoreParametersDataSource{}
	_ datasource.DataSourceWithConfigure = &DatastoreParametersDataSource{}
)

func NewDatastoreParametersDataSource() datasource.DataSource {
	return &DatastoreParametersDataSource{}
}

type DatastoreParametersDataSource struct {
	config clients.Config
}

type DatastoreParametersDataSourceModel struct {
	ID     types.String `tfsdk:"id"`
	Region types.String `tfsdk:"region"`

	DatastoreName      types.String              `tfsdk:"datastore_name"`
	DatastoreVersionID types.String              `tfsdk:"datastore_version_id"`
	Parameters         []DatastoreParameterModel `tfsdk:"parameters"`
}

type DatastoreParameterModel struct {
	Max             types.Float64 `tfsdk:"max"`
	Min             types.Float64 `tfsdk:"min"`
	Name            types.String  `tfsdk:"name"`
	RestartRequried types.Bool    `tfsdk:"restart_required"`
	Type            types.String  `tfsdk:"type"`
}

func (d *DatastoreParametersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_db_datastore_parameters"
}

func (d *DatastoreParametersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

			"datastore_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the data store.",
			},

			"datastore_version_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the version of the data store.",
			},

			"parameters": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"min": schema.Float64Attribute{
							Computed:    true,
							Description: "Minimum value of a configuration parameter.",
						},

						"max": schema.Float64Attribute{
							Computed:    true,
							Description: "Maximum value of a configuration parameter.",
						},

						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of a configuration parameter.",
						},

						"restart_required": schema.BoolAttribute{
							Computed:    true,
							Description: "This attribute indicates whether a restart required when a parameter is set.",
						},

						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of a configuration parameter.",
						},
					},
				},
				Description: "Configuration parameters supported for the datastore.",
			},
		},
		Description: "Use this data source to get configuration parameters supported for a VKCS datastore.",
	}
}

func (d *DatastoreParametersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *DatastoreParametersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DatastoreParametersDataSourceModel

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

	dsName, dsVersionID := data.DatastoreName.ValueString(), data.DatastoreVersionID.ValueString()
	ctx = tflog.SetField(ctx, "datastore_name", dsName)
	ctx = tflog.SetField(ctx, "datastore_version_id", dsVersionID)

	tflog.Debug(ctx, "Calling Databases API to list parameters for the datastore")

	params, err := datastores.ListParameters(client, dsName, dsVersionID).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling VKCS Databases API", err.Error())
		return
	}

	tflog.Debug(ctx, "Called Databases API to list parameters for the datastore", map[string]interface{}{"parameters": fmt.Sprintf("%#v", params)})

	data.ID = types.StringValue(fmt.Sprintf("%s/%s/params", dsName, dsVersionID))
	data.Region = types.StringValue(region)
	data.Parameters = flattenDatastoreParameters(params)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenDatastoreParameters(params []datastores.Param) (r []DatastoreParameterModel) {
	for _, p := range params {
		r = append(r, DatastoreParameterModel{
			Max:             types.Float64Value(p.MaxValue),
			Min:             types.Float64Value(p.MinValue),
			Name:            types.StringValue(p.Name),
			RestartRequried: types.BoolValue(p.RestartRequried),
			Type:            types.StringValue(p.Type),
		})
	}
	return
}
