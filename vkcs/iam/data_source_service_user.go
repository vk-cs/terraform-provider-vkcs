package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/iam/datasource_service_user"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/iam/serviceusers"
)

var (
	_ datasource.DataSource                     = (*serviceUserDataSource)(nil)
	_ datasource.DataSourceWithConfigure        = (*serviceUserDataSource)(nil)
	_ datasource.DataSourceWithConfigValidators = (*serviceUserDataSource)(nil)
)

func NewServiceUserDataSource() datasource.DataSource {
	return &serviceUserDataSource{}
}

type serviceUserDataSource struct {
	config clients.Config
}

func (d *serviceUserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_service_user"
}

func (d *serviceUserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_service_user.ServiceUserDataSourceSchema(ctx)
}

func (d *serviceUserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *serviceUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_service_user.ServiceUserModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	client, err := d.config.IAMServiceUsersV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating IAM Service Users API client", err.Error())
		return
	}

	var serviceUser *serviceusers.ServiceUser

	if id := data.Id.ValueString(); id != "" {
		ctx = tflog.SetField(ctx, "service_user_id", id)
		tflog.Trace(ctx, "Calling IAM Service Users API to get service user")

		var err error
		serviceUser, err = serviceusers.Get(client, id).Extract()
		if err != nil {
			resp.Diagnostics.AddError("Error calling IAM Service Users API to get service user", err.Error())
			return
		}

		tflog.Trace(ctx, "Called IAM Service Users API to get service user", map[string]any{"service_user": fmt.Sprintf("%#v", serviceUser)})
	} else {
		listOpts := serviceusers.ListOpts{
			Name: data.Name.ValueString(),
		}

		tflog.Trace(ctx, "Calling IAM Service Users API to list service users", map[string]any{"opts": fmt.Sprintf("%#v", listOpts)})

		allPages, err := serviceusers.List(client, listOpts).AllPages()
		if err != nil {
			resp.Diagnostics.AddError("Error calling IAM Service Users API to list service users", err.Error())
			return
		}

		tflog.Trace(ctx, "Called IAM Service Users API to list service users", map[string]any{"all_pages": fmt.Sprintf("%#v", allPages)})

		serviceUsers, err := serviceusers.ExtractServiceUsers(allPages)
		if err != nil {
			resp.Diagnostics.AddError("Error extracting service users", err.Error())
			return
		}

		tflog.Trace(ctx, "Extracted list of service users", map[string]any{"service_users": fmt.Sprintf("%#v", serviceUsers)})

		if len(serviceUsers) < 1 {
			resp.Diagnostics.AddError("Your query returned no results", "Please change your search criteria and try again")
			return
		}

		if len(serviceUsers) > 1 {
			resp.Diagnostics.AddError("Your query returned more than one result", "Please try a more specific search criteria")
			return
		}

		serviceUser = &serviceUsers[0]
	}

	resp.Diagnostics.Append(data.UpdateFromServiceUser(ctx, serviceUser)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *serviceUserDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}
