package images

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/framework/utils"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/framework/validators"
)

func dateFilters() []string {
	return []string{
		string(images.FilterEQ),
		string(images.FilterNEQ),
		string(images.FilterGT),
		string(images.FilterGTE),
		string(images.FilterLT),
		string(images.FilterLTE),
	}
}

var (
	_ datasource.DataSource              = &ImagesDataSource{}
	_ datasource.DataSourceWithConfigure = &ImagesDataSource{}
)

func NewImagesDataSource() datasource.DataSource {
	return &ImagesDataSource{}
}

type ImagesDataSource struct {
	config clients.Config
}

type ImagesDataSourceModel struct {
	ID     types.String `tfsdk:"id"`
	Region types.String `tfsdk:"region"`

	CreatedAt  types.String `tfsdk:"created_at"`
	Default    types.Bool   `tfsdk:"default"`
	Images     []ImageModel `tfsdk:"images"`
	Owner      types.String `tfsdk:"owner"`
	Properties types.Map    `tfsdk:"properties"`
	SizeMin    types.Int64  `tfsdk:"size_min"`
	SizeMax    types.Int64  `tfsdk:"size_max"`
	Tags       types.List   `tfsdk:"tags"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
	Visibility types.String `tfsdk:"visibility"`
}

type ImageModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Properties types.Map    `tfsdk:"properties"`
}

func (d *ImagesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_images_images"
}

func (d *ImagesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source",
			},

			"region": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The region in which to obtain the Images client. If omitted, the `region` argument of the provider is used.",
			},

			"created_at": schema.StringAttribute{
				Optional:   true,
				Validators: []validator.String{validators.DateFilter(dateFilters()...)},
				Description: fmt.Sprintf("Date filter to select images with created_at matching the specified criteria. "+
					"Value should be either RFC3339 formatted time or time filter in format `filter:time`, where "+
					"`filter` is one of [%s] and `time` is RFC3339 formatted time.", strings.Join(dateFilters(), ", ")),
			},

			"default": schema.BoolAttribute{
				Optional:    true,
				Description: "The flag used to filter images based on whether they are available for virtual machine creation.",
			},

			"images": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Images matching specified criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of an image.",
						},

						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of an image.",
						},

						"properties": schema.MapAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "Properties associated with an image.",
						},
					},
				},
			},

			"owner": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of the owner of images.",
			},

			"properties": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Search for images with specific properties.",
			},

			"size_min": schema.Int64Attribute{
				Optional:    true,
				Description: "The minimum size (in bytes) of images to return.",
			},

			"size_max": schema.Int64Attribute{
				Optional:    true,
				Description: "The maximum size (in bytes) of images to return.",
			},

			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Search for images with specific tags.",
			},

			"updated_at": schema.StringAttribute{
				Optional:   true,
				Validators: []validator.String{validators.DateFilter(dateFilters()...)},
				Description: fmt.Sprintf("Date filter to select images with updated_at matching the specified criteria. "+
					"Value should be either RFC3339 formatted time or time filter in format `filter:time`, where "+
					"`filter` is one of [%s] and `time` is RFC3339 formatted time.", strings.Join(dateFilters(), ", ")),
			},

			"visibility": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(images.ImageVisibilityPublic),
						string(images.ImageVisibilityPrivate),
						string(images.ImageVisibilityShared),
						string(images.ImageVisibilityCommunity),
					),
				},
				Description: "The visibility of images. Must be one of \"public\", \"private\", \"community\", or \"shared\".",
			},
		},
	}
}

func (d *ImagesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *ImagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ImagesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}

	client, err := d.config.ImageV2Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS image client", err.Error())
		return
	}

	visibility := data.Visibility.ValueString()

	var tags []string
	resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var properties map[string]string
	resp.Diagnostics.Append(data.Properties.ElementsAs(ctx, &properties, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	listOpts := ListOpts{
		Visibility: resourceImagesImageVisibilityFromString(visibility),
		Owner:      data.Owner.ValueString(),
		Status:     images.ImageStatusActive,
		SizeMin:    data.SizeMin.ValueInt64(),
		SizeMax:    data.SizeMax.ValueInt64(),
		Tags:       tags,
		Properties: properties,
	}
	if utils.IsKnown(data.CreatedAt) {
		listOpts.CreatedAtQuery = expandDateFilter(data.CreatedAt.ValueString())
	}
	if utils.IsKnown(data.UpdatedAt) {
		listOpts.UpdatedAtQuery = expandDateFilter(data.UpdatedAt.ValueString())
	}

	tflog.Debug(ctx, "Calling Images API to list images", map[string]interface{}{"list_opts": fmt.Sprintf("%#v", listOpts)})

	allPages, err := images.List(client, listOpts).AllPages()
	if err != nil {
		resp.Diagnostics.AddError("Error calling VKCS Images API", err.Error())
		return
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		resp.Diagnostics.AddError("Error processing VKSC Images API response", err.Error())
		return
	}

	tflog.Debug(ctx, "Called Images API to list images", map[string]interface{}{"images": fmt.Sprintf("%#v", allImages)})

	if data.Default.ValueBool() {
		allImages = filterImagesByDefault(allImages)
		tflog.Debug(ctx, "Filtered images by default flag", map[string]interface{}{"images": fmt.Sprintf("%#v", allImages)})
	}

	flattenedImages := flattenImages(ctx, allImages, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	sort.SliceStable(flattenedImages, func(i, j int) bool {
		return flattenedImages[i].Name.ValueString() < flattenedImages[j].Name.ValueString()
	})

	data.ID = types.StringValue(strconv.FormatInt(time.Now().Unix(), 10))
	data.Region = types.StringValue(region)
	data.Images = flattenedImages

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func expandDateFilter(date string) *images.ImageDateQuery {
	// error checks are not necessary, since they were validated by terraform validate functions
	var parts []string
	if regexp.MustCompile("^" + strings.Join(dateFilters(), "|") + ":").Match([]byte(date)) {
		parts = strings.SplitN(date, ":", 2)
	} else {
		parts = []string{date}
	}

	var parsedTime time.Time
	var filter *images.ImageDateQuery

	if len(parts) == 2 {
		parsedTime, _ = time.Parse(time.RFC3339, parts[1])
		filter = &images.ImageDateQuery{Date: parsedTime, Filter: images.ImageDateFilter(parts[0])}
	} else {
		parsedTime, _ = time.Parse(time.RFC3339, parts[0])
		filter = &images.ImageDateQuery{Date: parsedTime}
	}

	if parsedTime == (time.Time{}) {
		return nil
	}

	return filter
}

func flattenImages(ctx context.Context, images []images.Image, respDiags *diag.Diagnostics) []ImageModel {
	r := make([]ImageModel, 0)
	for _, im := range images {
		propsM := make(map[string]string, len(im.Properties))
		for k, v := range im.Properties {
			vStr, ok := v.(string)
			if !ok {
				continue
			}
			propsM[k] = vStr
		}
		properties, diags := types.MapValueFrom(ctx, types.StringType, propsM)
		respDiags.Append(diags...)
		if diags.HasError() {
			return nil
		}
		r = append(r, ImageModel{
			ID:         types.StringValue(im.ID),
			Name:       types.StringValue(im.Name),
			Properties: properties,
		})
	}
	return r
}
