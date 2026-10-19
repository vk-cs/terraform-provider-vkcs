package kms

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

var _ datasource.DataSourceWithConfigure = (*keysListDataSource)(nil)

func NewKeysListDataSource() datasource.DataSource {
	return &keysListDataSource{}
}

type keysListDataSource struct {
	config clients.Config
}

type keysListDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	KeysCount types.Int64  `tfsdk:"keys_count"`
	Keys      types.List   `tfsdk:"keys"`
}

func (r *keysListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = typeName(req.ProviderTypeName, KeysListEntityName)
}

func (r *keysListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData != nil {
		r.config = req.ProviderData.(clients.Config)
	}
}

func (r *keysListDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source containing list of KMS keys",
		Attributes: map[string]schema.Attribute{
			EntityFieldID: schema.StringAttribute{Computed: true, Description: "Identifier of the KMS keys list"},
			KeysListFieldKeysCount: schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Number of keys in list",
			},
			KeysListFieldKeys: schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of keys' names",
			},
		},
	}
}

func (r *keysListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data keysListDataSourceModel
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

	var keys []string
	err = waitForKMSOperation(ctx, func() (bool, error) {
		var err error
		keys, err = listKeys(kmsV1Client)
		if errutil.IsNotFound(err) {
			return false, nil
		}
		return err == nil, err
	})
	if err != nil {
		resp.Diagnostics.AddError("Error listing keys", err.Error())
		return
	}

	var builder strings.Builder
	for _, secret := range keys[:min(len(keys), maxItemsInListForIDBuild)] {
		_, err := builder.WriteString(secret)
		if err != nil {
			resp.Diagnostics.AddError("Error building id", err.Error())
			return
		}
	}
	hash := sha256.Sum256(fmt.Appendf(nil, "%s:%d", builder.String(), len(keys)))

	data.ID = types.StringValue(hex.EncodeToString(hash[:]))
	data.KeysCount = types.Int64Value(int64(len(keys)))
	if keys == nil {
		keys = []string{}
	}
	var listDiags diag.Diagnostics
	data.Keys, listDiags = types.ListValueFrom(ctx, types.StringType, keys)
	resp.Diagnostics.Append(listDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
