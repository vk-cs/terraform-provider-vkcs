package kms

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kms/v1/secrets"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

var _ datasource.DataSourceWithConfigure = (*secretDataSource)(nil)

func NewSecretDataSource() datasource.DataSource {
	return &secretDataSource{}
}

type secretDataSource struct {
	config clients.Config
}

type secretDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	Path             types.String `tfsdk:"path"`
	DataJSON         types.String `tfsdk:"data_json"`
	CreatedTime      types.String `tfsdk:"created_time"`
	Version          types.Int64  `tfsdk:"version"`
	DeleteProtection types.Bool   `tfsdk:"delete_protection"`
}

func (r *secretDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kms_secret"
}

func (r *secretDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData != nil {
		r.config = req.ProviderData.(clients.Config)
	}
}

func (r *secretDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source containing a KMS secret",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Identifier of the KMS secret"},
			"path": schema.StringAttribute{
				Required:    true,
				Description: "Identifier (path) of the secret within the KMS secret store.",
			},
			"data_json": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "Secret data as a JSON object, represented as a string and populated from KMS when the data source is read. Use `jsondecode` to access its fields.",
			},
			"created_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time returned in the secret metadata, formatted as a Go time string, for example `2026-09-01 10:00:00 +0000 UTC`.",
			},
			"version": schema.Int64Attribute{
				Computed:    true,
				Description: "Current version of the secret returned in its metadata. This field does not select a version to read.",
			},
			"delete_protection": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the secret is protected from deletion.",
			},
		},
	}
}

func (r *secretDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data secretDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config := r.config
	kmsV1Client, err := config.KMSV1Client(config.GetRegion())
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS KMS client", err.Error())
		return
	}

	path := data.Path.ValueString()

	var secret secrets.GetSecret
	err = waitForKMSOperation(ctx, func() (bool, error) {
		var err error
		secret, err = getSecret(kmsV1Client, path)
		if errutil.IsNotFound(err) {
			return false, nil
		}
		return err == nil, err
	})
	if err != nil {
		resp.Diagnostics.AddError("Error getting secret", err.Error())
		return
	}

	var deleteProtection bool
	err = waitForKMSOperation(ctx, func() (bool, error) {
		var err error
		deleteProtection, err = getSecretDeleteProtection(kmsV1Client, path)
		if errutil.IsNotFound(err) {
			return false, nil
		}
		return err == nil, err
	})
	if err != nil {
		resp.Diagnostics.AddError("Error getting secret delete protection", err.Error())
		return
	}

	data.ID = types.StringValue(path)

	secretString, err := json.Marshal(secret.Data.Data)
	if err != nil {
		resp.Diagnostics.AddError("Error marshalling data_json", err.Error())
		return
	}

	data.DataJSON = types.StringValue(string(secretString))

	data.CreatedTime = types.StringValue(secret.Data.Metadata.CreatedTime.String())

	data.Version = types.Int64Value(int64(secret.Data.Metadata.Version))

	data.DeleteProtection = types.BoolValue(deleteProtection)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
