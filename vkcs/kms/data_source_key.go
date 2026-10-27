package kms

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kms/v1/keys"
)

var _ datasource.DataSourceWithConfigure = (*keyDataSource)(nil)

func NewKeyDataSource() datasource.DataSource {
	return &keyDataSource{}
}

type keyDataSource struct {
	config clients.Config
}

type keyDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	DeletionAllowed types.Bool   `tfsdk:"deletion_allowed"`
}

func (r *keyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kms_key"
}

func (r *keyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData != nil {
		r.config = req.ProviderData.(clients.Config)
	}
}

func (r *keyDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source containing a KMS key",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Description: "Identifier of the KMS key"},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the KMS key to read.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "Type of the key: `aes128-gcm96`, `aes256-gcm96`, `chacha20-poly1305`, or `xchacha20-poly1305`.",
			},
			"deletion_allowed": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether KMS allows the key to be deleted.",
			},
		},
	}
}

func (r *keyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data keyDataSourceModel
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

	name := data.Name.ValueString()

	result, err := kmsStateChangeConf(keyStateRefreshFunc(kmsV1Client, name, nil), kmsStateReady).WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error getting key", err.Error())
		return
	}
	key, ok := result.(keys.GetKey)
	if !ok {
		resp.Diagnostics.AddError("Unexpected KMS response type", fmt.Sprintf("Expected keys.GetKey, got %T", result))
		return
	}

	data.ID = types.StringValue(name)

	data.Name = types.StringValue(key.Data.Name)

	data.Type = types.StringValue(key.Data.Type)

	data.DeletionAllowed = types.BoolValue(key.Data.DeletionAllowed)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
