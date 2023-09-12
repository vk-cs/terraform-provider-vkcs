package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/datastores"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

// Ensure the implementation satisfies the desired interfaces.
var _ datasource.DataSource = &DatastoreCapabilitiesDataSource{}

func NewDatastoreCapabilitiesDataSource() datasource.DataSource {
	return &DatastoreCapabilitiesDataSource{}
}

type DatastoreCapabilitiesDataSource struct {
	config clients.Config
}

type DatastoreCapabilitiesDataSourceModel struct {
	ID                 types.String                   `tfsdk:"id"`
	DatastoreName      types.String                   `tfsdk:"datastore_name"`
	DatastoreVersionID types.String                   `tfsdk:"datastore_version_id"`
	Capabilities       []DatastoreCapabilityItemModel `tfsdk:"capabilities"`
	Region             types.String                   `tfsdk:"region"`
}

type DatastoreCapabilityItemModel struct {
	Name                   types.String                     `tfsdk:"name"`
	Description            types.String                     `tfsdk:"description"`
	Params                 []DatastoreCapabilityParamsModel `tfsdk:"params"`
	ShouldBeOnMaster       types.Bool                       `tfsdk:"should_be_on_master"`
	AllowMajorUpgrade      types.Bool                       `tfsdk:"allow_major_upgrade"`
	AllowUpgradeFromBackup types.Bool                       `tfsdk:"allow_upgrade_from_backup"`
}

type DatastoreCapabilityParamsModel struct {
	Name         types.String  `tfsdk:"name"`
	Required     types.Bool    `tfsdk:"required"`
	Type         types.String  `tfsdk:"type"`
	ElementType  types.String  `tfsdk:"element_type"`
	EnumValues   types.List    `tfsdk:"enum_values"`
	DefaultValue types.String  `tfsdk:"default_value"`
	Min          types.Float64 `tfsdk:"min"`
	Max          types.Float64 `tfsdk:"max"`
	Regex        types.String  `tfsdk:"regex"`
	Masked       types.Bool    `tfsdk:"masked"`
}

func (d *DatastoreCapabilitiesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_db_datastore_capabilities"
}

func (d *DatastoreCapabilitiesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource",
			},

			"region": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The `region` to fetch availability zones from, defaults to the provider's `region`.",
			},

			"datastore_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the data store.",
			},

			"datastore_version_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the version of the data store.",
			},

			"capabilities": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of data store capability.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Description of data store capability.",
						},
						"params": schema.ListNestedAttribute{
							Computed:     true,
							NestedObject: DatastoreCapabilitiesParamSchema(),
						},
						"should_be_on_master": schema.BoolAttribute{
							Computed:    true,
							Description: "This attribute indicates whether a capability applies only to the master node.",
						},
						"allow_major_upgrade": schema.BoolAttribute{
							Computed:    true,
							Description: "This attribute indicates whether a capability can be applied in the next major version of data store.",
						},
						"allow_upgrade_from_backup": schema.BoolAttribute{
							Computed:    true,
							Description: "This attribute indicates whether a capability can be applied to upgrade from backup.",
						},
					},
				},
				Description: "Capabilities of the datastore.",
			},
		},
		Description: "Use this data source to get capabilities supported for a VKCS datastore.",
	}
}

func DatastoreCapabilitiesParamSchema() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of a parameter.",
			},
			"required": schema.BoolAttribute{
				Computed:    true,
				Description: "Required indicates whether a parameter value must be set.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "Type of value for a parameter.",
			},
			"element_type": schema.StringAttribute{
				Computed:    true,
				Description: "Type of element value for a parameter of `list` type.",
			},
			"enum_values": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Supported values for a parameter.",
			},
			"default_value": schema.StringAttribute{
				Computed:    true,
				Description: "Default value for a parameter.",
			},
			"min": schema.Float64Attribute{
				Computed:    true,
				Description: "Minimum value for a parameter.",
			},
			"max": schema.Float64Attribute{
				Computed:    true,
				Description: "Maximum value for a parameter.",
			},
			"regex": schema.StringAttribute{
				Computed:    true,
				Description: "Regular expression that a parameter value must match.",
			},
			"masked": schema.BoolAttribute{
				Computed:    true,
				Description: "Masked indicates whether a parameter value must be a boolean mask.",
			},
		},
	}
}

func (d *DatastoreCapabilitiesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *DatastoreCapabilitiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DatastoreCapabilitiesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dsName := data.DatastoreName.ValueString()
	dsVersionID := data.DatastoreVersionID.ValueString()
	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}

	dbClient, err := d.config.DatabaseV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS database client", err.Error())
		return
	}

	capabilities, err := datastores.ListCapabilities(dbClient, dsName, dsVersionID).Extract()
	if err != nil {
		checkDeleted := util.CheckDeletedDatasource(ctx, resp, err)
		if checkDeleted != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_db_datastore_capabilities", checkDeleted.Error())
		}

		return
	}
	flattenedCapabilities, diags := flattenDatastoreCapabilities(ctx, capabilities)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Capabilities = flattenedCapabilities

	data.ID = types.StringValue(fmt.Sprintf("%s/%s/capabilities", dsName, dsVersionID))
	data.Region = types.StringValue(region)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenDatastoreCapabilities(ctx context.Context, capabilities []datastores.Capability) (r []DatastoreCapabilityItemModel, diags diag.Diagnostics) {
	for _, c := range capabilities {
		params, diags := flattenDatastoreCapabilityParams(ctx, c.Params)
		if diags.HasError() {
			return nil, diags
		}

		r = append(r, DatastoreCapabilityItemModel{
			Name:                   types.StringValue(c.Name),
			Description:            types.StringValue(c.Description),
			Params:                 params,
			ShouldBeOnMaster:       types.BoolValue(c.ShouldBeOnMaster),
			AllowMajorUpgrade:      types.BoolValue(c.AllowMajorUpgrade),
			AllowUpgradeFromBackup: types.BoolValue(c.AllowUpgradeFromBackup),
		})
	}
	return
}

func flattenDatastoreCapabilityParams(ctx context.Context, params map[string]*datastores.CapabilityParam) (r []DatastoreCapabilityParamsModel, diags diag.Diagnostics) {
	for name, p := range params {
		var defaultValue string
		switch v := p.DefaultValue.(type) {
		case string:
			defaultValue = v
		case float64:
			defaultValue = strconv.FormatFloat(p.DefaultValue.(float64), 'f', -1, 64)
		}

		enumValues, diags := types.ListValueFrom(ctx, types.StringType, p.EnumValues)
		if diags.HasError() {
			return nil, diags
		}

		r = append(r, DatastoreCapabilityParamsModel{
			Name:         types.StringValue(name),
			Required:     types.BoolValue(p.Required),
			Type:         types.StringValue(p.Type),
			ElementType:  types.StringValue(p.ElementType),
			EnumValues:   enumValues,
			DefaultValue: types.StringValue(defaultValue),
			Min:          types.Float64Value(p.MinValue),
			Max:          types.Float64Value(p.MaxValue),
			Regex:        types.StringValue(p.Regex),
			Masked:       types.BoolValue(p.Masked),
		})
	}
	return
}
