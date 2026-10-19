package kms

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

var _ datasource.DataSourceWithConfigure = (*keyEncryptDataSource)(nil)

func NewKeyEncryptDataSource() datasource.DataSource {
	return &keyEncryptDataSource{}
}

type keyEncryptDataSource struct {
	config clients.Config
}

type keyEncryptDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	Key        types.String `tfsdk:"key"`
	Ciphertext types.String `tfsdk:"ciphertext"`
	Plaintext  types.String `tfsdk:"plaintext"`
}

func (r *keyEncryptDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = typeName(req.ProviderTypeName, KeyEncryptEntityName)
}

func (r *keyEncryptDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData != nil {
		r.config = req.ProviderData.(clients.Config)
	}
}

func (r *keyEncryptDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source containing encrypt result by KMS key",
		Attributes: map[string]schema.Attribute{
			EntityFieldID: schema.StringAttribute{Computed: true, Description: "Identifier of the KMS key encrypt"},
			KeyCryptFieldKey: schema.StringAttribute{
				Required:    true,
				Description: "Key used to encrypt data",
			},
			KeyCryptFieldCiphertext: schema.StringAttribute{
				Computed:    true,
				Description: "Resulting encrypted data",
			},
			KeyCryptFieldPlaintext: schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Data to encrypt encoded in base64",
			},
		},
	}
}

func (r *keyEncryptDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data keyEncryptDataSourceModel
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

	plaintext := data.Plaintext.ValueString()

	var ciphertext string
	err = waitForKMSOperation(ctx, func() (bool, error) {
		var err error
		ciphertext, err = encrypt(kmsV1Client, key, plaintext)
		if errutil.IsNotFound(err) {
			return false, nil
		}
		return err == nil, err
	})
	if err != nil {
		resp.Diagnostics.AddError("Error encrypting", err.Error())
		return
	}

	hash := sha256.Sum256([]byte(ciphertext))
	data.ID = types.StringValue(hex.EncodeToString(hash[:]))

	data.Ciphertext = types.StringValue(ciphertext)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
