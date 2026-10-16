package kms

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
)

var _ datasource.DataSourceWithConfigure = (*keyDecryptDataSource)(nil)

func NewKeyDecryptDataSource() datasource.DataSource {
	return &keyDecryptDataSource{}
}

type keyDecryptDataSource struct {
	config clients.Config
}

type keyDecryptDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	Key        types.String `tfsdk:"key"`
	Ciphertext types.String `tfsdk:"ciphertext"`
	Plaintext  types.String `tfsdk:"plaintext"`
}

func (r *keyDecryptDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = typeName(req.ProviderTypeName, KeyDecryptEntityName)
}

func (r *keyDecryptDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData != nil {
		r.config = req.ProviderData.(clients.Config)
	}
}

func (r *keyDecryptDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source containing decrypt result by KMS key",
		Attributes: map[string]schema.Attribute{
			EntityFieldID: schema.StringAttribute{Computed: true, Description: "Identifier of the KMS key decrypt"},
			KeyCryptFieldKey: schema.StringAttribute{
				Required:    true,
				Description: "Key used to decrypt data",
			},
			KeyCryptFieldCiphertext: schema.StringAttribute{
				Required:    true,
				Description: "Encrypted data",
			},
			KeyCryptFieldPlaintext: schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "Resulting decrypted data encoded in base64",
			},
		},
	}
}

func (r *keyDecryptDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data keyDecryptDataSourceModel
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

	key := data.Key.ValueString()

	ciphertext := data.Ciphertext.ValueString()

	plaintext, err := decrypt(kmsV1Client, key, ciphertext)
	if err != nil {
		resp.Diagnostics.AddError("Error decrypting", err.Error())
		return
	}

	hash := sha256.Sum256([]byte(plaintext))
	data.ID = types.StringValue(hex.EncodeToString(hash[:]))

	data.Plaintext = types.StringValue(plaintext)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
