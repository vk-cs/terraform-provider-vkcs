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
)

var _ datasource.DataSourceWithConfigure = (*secretsListDataSource)(nil)

func NewSecretsListDataSource() datasource.DataSource {
	return &secretsListDataSource{}
}

type secretsListDataSource struct {
	config clients.Config
}

type secretsListDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	SecretsCount types.Int64  `tfsdk:"secrets_count"`
	Secrets      types.List   `tfsdk:"secrets"`
}

func (r *secretsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = typeName(req.ProviderTypeName, SecretsListEntityName)
}

func (r *secretsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData != nil {
		r.config = req.ProviderData.(clients.Config)
	}
}

func (r *secretsListDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source containing list of KMS secrets",
		Attributes: map[string]schema.Attribute{
			EntityFieldID: schema.StringAttribute{Computed: true, Description: "Identifier of the KMS secrets list"},
			SecretsListFieldSecretsCount: schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Number of secrets in list",
			},
			SecretsListFieldSecrets: schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of keys' paths",
			},
		},
	}
}

func (r *secretsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data secretsListDataSourceModel
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

	var secrets []string
	err = waitForKMSOperation(ctx, func() (bool, error) {
		var err error
		secrets, err = listSecrets(kmsV1Client)
		return err == nil, err
	})
	if err != nil {
		resp.Diagnostics.AddError("Error listing secrets", err.Error())
		return
	}

	var builder strings.Builder
	for _, secret := range secrets[:min(len(secrets), maxItemsInListForIDBuild)] {
		_, err := builder.WriteString(secret)
		if err != nil {
			resp.Diagnostics.AddError("Error building id", err.Error())
			return
		}
	}
	hash := sha256.Sum256(fmt.Appendf(nil, "%s:%d", builder.String(), len(secrets)))

	data.ID = types.StringValue(hex.EncodeToString(hash[:]))
	data.SecretsCount = types.Int64Value(int64(len(secrets)))
	if secrets == nil {
		secrets = []string{}
	}
	var listDiags diag.Diagnostics
	data.Secrets, listDiags = types.ListValueFrom(ctx, types.StringType, secrets)
	resp.Diagnostics.Append(listDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
