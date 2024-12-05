package db

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/framework/utils"
	configgroups "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/config_groups"
)

var (
	_ datasource.DataSource              = &ConfigGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &ConfigGroupDataSource{}
)

func NewConfigGroupDataSource() datasource.DataSource {
	return &ConfigGroupDataSource{}
}

type ConfigGroupDataSource struct {
	config clients.Config
}

type ConfigGroupDataSourceModel struct {
	ID     types.String `tfsdk:"id"`
	Region types.String `tfsdk:"region"`

	ConfigGroupID types.String                `tfsdk:"config_group_id"`
	Created       types.String                `tfsdk:"created"`
	Datastore     []ConfigGroupDatastoreModel `tfsdk:"datastore"`
	Description   types.String                `tfsdk:"description"`
	Name          types.String                `tfsdk:"name"`
	Updated       types.String                `tfsdk:"updated"`
	Values        types.Map                   `tfsdk:"values"`
}

type ConfigGroupDatastoreModel struct {
	Type    types.String `tfsdk:"type"`
	Version types.String `tfsdk:"version"`
}

func (d *ConfigGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_db_config_group"
}

func (d *ConfigGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The UUID of the config_group.",
			},

			"region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The region in which to obtain the service client. If omitted, the `region` argument of the provider is used.",
			},

			"config_group_id": schema.StringAttribute{
				Optional:           true,
				Computed:           true,
				Description:        "The UUID of the config_group.",
				DeprecationMessage: "This argument is deprecated, please, use the `id` attribute instead.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("id"),
					),
				},
			},

			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp of config group's creation.",
			},

			"datastore": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required:    true,
							Description: "Type of the datastore.",
						},

						"version": schema.StringAttribute{
							Required:    true,
							Description: "Version of the datastore.",
						},
					},
				},
				Computed:    true,
				Description: "Object that represents datastore of backup",
			},

			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The description of the config group.",
			},

			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the config group.",
			},

			"updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp of config group's last update.",
			},

			"values": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Map of configuration parameters in format \"key\": \"value\".",
			},
		},
		Description: "Use this data source to get the information on a db config group resource.",
	}
}

func (d *ConfigGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *ConfigGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConfigGroupDataSourceModel

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

	configGroupID := utils.GetFirstNotEmptyValue(data.ID, data.ConfigGroupID)
	ctx = tflog.SetField(ctx, "id", configGroupID)

	tflog.Debug(ctx, "Calling Databases API to read the config group")

	configGroup, err := configgroups.Get(client, configGroupID).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling VKCS Databases API", err.Error())
		return
	}

	tflog.Debug(ctx, "Called Databases API to read the config group", map[string]any{"config_group": fmt.Sprintf("%#v", configGroup)})

	data.ID = types.StringValue(configGroupID)
	data.ConfigGroupID = data.ID
	data.Region = types.StringValue(region)
	data.Created = types.StringValue(configGroup.Created)
	data.Datastore = flattenConfigGroupDatastore(configGroup)
	data.Description = types.StringValue(configGroup.Description)
	data.Name = types.StringValue(configGroup.Name)
	data.Updated = types.StringValue(configGroup.Updated)
	data.Values = flattenConfigGroupValues(ctx, configGroup.Values, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenConfigGroupDatastore(cg *configgroups.ConfigGroupResp) []ConfigGroupDatastoreModel {
	if cg == nil {
		return nil
	}
	return []ConfigGroupDatastoreModel{
		{
			Type:    types.StringValue(cg.DatastoreName),
			Version: types.StringValue(cg.DatastoreVersionName),
		},
	}
}

func flattenConfigGroupValues(ctx context.Context, v map[string]interface{}, respDiags *diag.Diagnostics) types.Map {
	rawValues := make(map[string]string, len(v))
	for name, value := range v {
		rawValues[name] = fmt.Sprintf("%v", value)
	}
	values, diags := types.MapValueFrom(ctx, types.StringType, rawValues)
	respDiags.Append(diags...)
	return values
}
