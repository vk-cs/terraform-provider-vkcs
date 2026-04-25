package baremetal

import (
	"context"
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/baremetal/v1/images"
)

var _ datasource.DataSource = &OSDataSource{}
var _ datasource.DataSourceWithConfigure = &OSDataSource{}

func NewOSDataSource() datasource.DataSource {
	return &OSDataSource{}
}

type OSDataSource struct {
	config clients.Config
}

type OSDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Region   types.String `tfsdk:"region"`
	Version  types.String `tfsdk:"version"`
	RaidType types.String `tfsdk:"raid_type"`
}

func (d *OSDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_baremetal_os"
}

func (d *OSDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The UUID of the OS.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("name"),
						path.MatchRoot("version"),
						path.MatchRoot("raid_type"),
					}...),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the OS.",
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Description: "The region to fetch the bare metal OS from, defaults to the provider's region.",
			},
			"version": schema.StringAttribute{
				Optional:    true,
				Description: "The version of the OS.",
			},
			"raid_type": schema.StringAttribute{
				Optional:    true,
				Description: "The raid type of the OS.",
			},
		},
		Description: "Use this data source to get information about a VKCS baremetal OS.",
	}
}

func (d *OSDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.config = req.ProviderData.(clients.Config)
}

func (d *OSDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OSDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}

	client, err := d.config.BareMetalV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS baremetal API client", err.Error())
		return
	}

	image, err := getImage(client, data)
	if err != nil {
		resp.Diagnostics.AddError("Error getting image", err.Error())
		return
	}

	data.ID = types.StringValue(image.ImageId)
	data.Region = types.StringValue(region)
	data.Name = types.StringValue(image.OsType)
	data.Version = types.StringValue(image.OsVersion)
	data.RaidType = types.StringValue(image.RaidType)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func getImage(client *gophercloud.ServiceClient, data OSDataSourceModel) (*images.Image, error) {
	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		return images.Get(client, data.ID.ValueString()).Extract()
	}

	imagePages, err := images.List(client).AllPages()
	if err != nil {
		return nil, err
	}

	imgs, err := images.ExtractImages(imagePages)
	if err != nil {
		return nil, err
	}

	filterImgs := filterImages(data, imgs)
	if len(filterImgs) == 0 {
		return nil, fmt.Errorf("no images found for OS")
	}

	if len(filterImgs) > 1 {
		return nil, fmt.Errorf("multiple images found for OS")
	}

	return &filterImgs[0], nil
}

func filterImages(data OSDataSourceModel, items []images.Image) (r []images.Image) {
	for _, item := range items {
		if imageMatches(data, item) {
			r = append(r, item)
		}
	}

	return
}

func imageMatches(data OSDataSourceModel, item images.Image) bool {
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		if !strings.EqualFold(item.OsType, data.Name.ValueString()) {
			return false
		}
	}
	if !data.Version.IsNull() && !data.Version.IsUnknown() {
		if !strings.EqualFold(item.OsVersion, data.Version.ValueString()) {
			return false
		}
	}
	if !data.RaidType.IsNull() && !data.RaidType.IsUnknown() {
		if !strings.EqualFold(item.RaidType, data.RaidType.ValueString()) {
			return false
		}
	}
	return true
}
