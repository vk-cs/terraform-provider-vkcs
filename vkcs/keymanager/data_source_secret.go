package keymanager

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/acls"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/framework/validators"
)

func dateFilters() []string {
	return []string{
		string(secrets.DateFilterGT),
		string(secrets.DateFilterGTE),
		string(secrets.DateFilterLT),
		string(secrets.DateFilterLTE),
	}
}

var (
	_ datasource.DataSource              = &SecretDataSource{}
	_ datasource.DataSourceWithConfigure = &SecretDataSource{}
)

func NewSecretDataSource() datasource.DataSource {
	return &SecretDataSource{}
}

type SecretDataSource struct {
	config clients.Config
}

type SecretDataSourceModel struct {
	ID     types.String `tfsdk:"id"`
	Region types.String `tfsdk:"region"`

	ACL                    []SecretDataSourceACLModel `tfsdk:"acl"`
	ACLOnly                types.Bool                 `tfsdk:"acl_only"`
	Algorithm              types.String               `tfsdk:"algorithm"`
	BitLength              types.Int64                `tfsdk:"bit_length"`
	ContentTypes           types.Map                  `tfsdk:"content_types"`
	CreatedAt              types.String               `tfsdk:"created_at"`
	CreatedAtFilter        types.String               `tfsdk:"created_at_filter"`
	CreatorID              types.String               `tfsdk:"creator_id"`
	Expiration             types.String               `tfsdk:"expiration"`
	ExpirationFilter       types.String               `tfsdk:"expiration_filter"`
	Metadata               types.Map                  `tfsdk:"metadata"`
	Mode                   types.String               `tfsdk:"mode"`
	Name                   types.String               `tfsdk:"name"`
	Payload                types.String               `tfsdk:"payload"`
	PayloadContentEncoding types.String               `tfsdk:"payload_content_encoding"`
	PayloadContentType     types.String               `tfsdk:"payload_content_type"`
	SecretRef              types.String               `tfsdk:"secret_ref"`
	SecretType             types.String               `tfsdk:"secret_type"`
	Status                 types.String               `tfsdk:"status"`
	UpdatedAt              types.String               `tfsdk:"updated_at"`
	UpdatedAtFilter        types.String               `tfsdk:"updated_at_filter"`
}

type SecretDataSourceACLModel struct {
	Read []SecretDataSourceACLOperationModel `tfsdk:"read"`
}

type SecretDataSourceACLOperationModel struct {
	CreatedAt     types.String `tfsdk:"created_at"`
	ProjectAccess types.Bool   `tfsdk:"project_access"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
	Users         types.Set    `tfsdk:"users"`
}

func (d *SecretDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_keymanager_secret"
}

func (d *SecretDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource.",
			},

			"region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The region in which to obtain the service client. If omitted, the `region` argument of the provider is used.",
			},

			"acl": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: d.ACLSchemas(),
				},
				Computed:    true,
				Description: "The list of ACLs assigned to a secret.",
			},

			"acl_only": schema.BoolAttribute{
				Optional:    true,
				Description: "Select the Secret with an ACL that contains the user. Project scope is ignored. Defaults to `false`.",
			},

			"algorithm": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The Secret algorithm.",
			},

			"bit_length": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The Secret bit length.",
			},

			"content_types": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "The map of the content types, assigned on the secret.",
			},

			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "The date the secret was created.",
			},

			"created_at_filter": schema.StringAttribute{
				Optional:    true,
				Description: "Date filter to select the Secret with created matching the specified criteria. See Date Filters below for more detail.",
				Validators:  []validator.String{validators.DateFilter(dateFilters()...)},
			},

			"creator_id": schema.StringAttribute{
				Computed:    true,
				Description: "The creator of the secret.",
			},

			"expiration": schema.StringAttribute{
				Computed:    true,
				Description: "The date the secret will expire.",
			},

			"expiration_filter": schema.StringAttribute{
				Optional:    true,
				Description: "Date filter to select the Secret with expiration matching the specified criteria. See Date Filters below for more detail.",
				Validators:  []validator.String{validators.DateFilter(dateFilters()...)},
			},

			"metadata": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "The map of metadata, assigned on the secret, which has been explicitly and implicitly added.",
			},

			"mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The Secret mode.",
			},

			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The Secret name.",
			},

			"payload": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The secret payload.",
			},

			"payload_content_encoding": schema.StringAttribute{
				Computed:    true,
				Description: "The Secret encoding.",
			},

			"payload_content_type": schema.StringAttribute{
				Computed:    true,
				Description: "The Secret content type.",
			},

			"secret_ref": schema.StringAttribute{
				Computed:    true,
				Description: "The secret reference / where to find the secret.",
			},

			"secret_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The Secret type. For more information see [Secret types](https://docs.openstack.org/barbican/latest/api/reference/secret_types.html).",
				Validators: []validator.String{
					stringvalidator.OneOf("symmetric", "public", "private", "passphrase", "certificate", "opaque"),
				},
			},

			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the secret.",
			},

			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "The date the secret was last updated.",
			},

			"updated_at_filter": schema.StringAttribute{
				Optional:    true,
				Description: "Date filter to select the Secret with updated matching the specified criteria. See Date Filters below for more detail.",
				Validators:  []validator.String{validators.DateFilter(dateFilters()...)},
			},
		},
		Description: "Use this data source to get the ID and the payload of an available Key secret\n\n~> **Important Security Notice** The payload of this data source will be stored *unencrypted* in your Terraform state file. **Use of this resource for production deployments is *not* recommended**. [Read more about sensitive data in state](https://www.terraform.io/docs/language/state/sensitive-data.html).",
	}
}

func (d *SecretDataSource) ACLSchemas() map[string]schema.Attribute {
	supportedACLOps := getSupportedACLOperations()
	aclSchemas := make(map[string]schema.Attribute, len(supportedACLOps))

	for _, op := range supportedACLOps {
		aclSchemas[op] = schema.ListNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"created_at": schema.StringAttribute{
						Computed:    true,
						Description: "The date the container ACL was created.",
					},

					"project_access": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether the container is accessible project wide. Defaults to `true`.",
					},

					"updated_at": schema.StringAttribute{
						Computed:    true,
						Description: "The date the container ACL was last updated.",
					},

					"users": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "The list of user IDs, which are allowed to access the container, when `project_access` is set to `false`.",
					},
				},
			},
			Computed:    true,
			Description: fmt.Sprintf("Block that describes %s operation.", op),
			Validators: []validator.List{
				listvalidator.SizeAtMost(1),
			},
		}
	}

	return aclSchemas
}

func (d *SecretDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *SecretDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SecretDataSourceModel
	var diags diag.Diagnostics

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Consider adding "region" attribute if it is not present.
	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}

	client, err := d.config.KeyManagerV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS Key Manager API client", err.Error())
		return
	}

	aclOnly := data.ACLOnly.ValueBool()
	listOpts := secrets.ListOpts{
		Name:            data.Name.ValueString(),
		Bits:            int(data.BitLength.ValueInt64()),
		Alg:             data.Algorithm.ValueString(),
		Mode:            data.Mode.ValueString(),
		SecretType:      secrets.SecretType(data.SecretType.ValueString()),
		ACLOnly:         &aclOnly,
		CreatedQuery:    expandDateFilter(data.CreatedAtFilter.ValueString()),
		UpdatedQuery:    expandDateFilter(data.UpdatedAtFilter.ValueString()),
		ExpirationQuery: expandDateFilter(data.ExpirationFilter.ValueString()),
	}

	tflog.Debug(ctx, "Calling Key Manager API to list secrets", map[string]interface{}{"list_opts": fmt.Sprintf("%#v", listOpts)})

	allPages, err := secrets.List(client, &listOpts).AllPages()
	if err != nil {
		resp.Diagnostics.AddError("Error calling VKCS Key Manager API", err.Error())
		return
	}

	allSecrets, err := secrets.ExtractSecrets(allPages)
	if err != nil {
		resp.Diagnostics.AddError("Error processing VKCS Key Manager API response", err.Error())
		return
	}

	tflog.Debug(ctx, "Called Key Manager API to list secrets", map[string]interface{}{"all_secrets_len": len(allSecrets)})

	if len(allSecrets) < 1 {
		resp.Diagnostics.AddError("Your query returned no vkcs_keymanager_secret results",
			"Please change your search criteria and try again")
		return
	}

	if len(allSecrets) > 1 {
		resp.Diagnostics.AddError("Your query returned more than one result",
			"Please try a more specific search criteria")
		return
	}

	secret := allSecrets[0]
	id := GetUUIDFromSecretRef(secret.SecretRef)
	ctx = tflog.SetField(ctx, "id", id)

	tflog.Debug(ctx, "Retrieved the secret", map[string]interface{}{"secret": fmt.Sprintf("%#v", secret)})
	tflog.Debug(ctx, "Calling Key Manager API to get the payload")

	payload, err := secrets.GetPayload(client, id, nil).Extract()
	if err != nil {
		tflog.Debug(ctx, "Error calling Key Manager API to get the payload", map[string]interface{}{"error": err.Error()})
	} else {
		tflog.Debug(ctx, "Called Key Manager API to get the payload", map[string]interface{}{"payload": string(payload)})
	}

	tflog.Debug(ctx, "Calling Key Manager API to get the metadata")

	metadata, err := secrets.GetMetadata(client, id).Extract()
	if err != nil {
		tflog.Debug(ctx, "Error calling Key Manager API to get the metadata", map[string]interface{}{"error": err.Error()})
	} else {
		tflog.Debug(ctx, "Called Key Manager API to get the metadata", map[string]interface{}{"metadata": fmt.Sprintf("%#v", metadata)})
	}

	tflog.Debug(ctx, "Calling Key Manager API to get the ACL")

	acl, err := acls.GetSecretACL(client, id).Extract()
	if err != nil {
		tflog.Debug(ctx, "Error calling Key Manager API to get the ACL", map[string]interface{}{"error": err.Error()})
	} else {
		tflog.Debug(ctx, "Called Key Manager API to get the ACL", map[string]interface{}{"acl": fmt.Sprintf("%#v", acl)})
	}

	data.ID = types.StringValue(id)
	data.Region = types.StringValue(region)
	data.ACL = flattenSecretDataSourceACL(ctx, acl, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Algorithm = types.StringValue(secret.Algorithm)
	data.BitLength = types.Int64Value(int64(secret.BitLength))
	data.ContentTypes, diags = types.MapValueFrom(ctx, types.StringType, secret.ContentTypes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.CreatedAt = types.StringValue(secret.Created.Format(time.RFC3339))
	data.CreatorID = types.StringValue(secret.CreatorID)
	if secret.Expiration != (time.Time{}) {
		data.Expiration = types.StringValue(secret.Expiration.Format(time.RFC3339))
	} else {
		data.Expiration = types.StringValue("")
	}
	data.Metadata, diags = types.MapValueFrom(ctx, types.StringType, metadata)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Mode = types.StringValue(secret.Mode)
	data.Name = types.StringValue(secret.Name)
	data.Payload = types.StringValue(string(payload))
	data.PayloadContentType = types.StringValue(secret.ContentTypes["default"])
	data.SecretType = types.StringValue(secret.SecretType)
	data.SecretRef = types.StringValue(secret.SecretRef)
	data.Status = types.StringValue(secret.Status)
	data.UpdatedAt = types.StringValue(secret.Updated.Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenSecretDataSourceACL(ctx context.Context, aclsPtr *acls.ACL, respDiags *diag.Diagnostics) []SecretDataSourceACLModel {
	if aclsPtr == nil {
		return nil
	}

	acls := *aclsPtr
	m := make(map[string][]SecretDataSourceACLOperationModel)

	for _, aclOp := range getSupportedACLOperations() {
		if v, ok := acls[aclOp]; ok {
			acl := flattenSecretDataSourceACLOperation(ctx, v, respDiags)
			if respDiags.HasError() {
				return nil
			}
			m[aclOp] = acl
		}
	}

	return []SecretDataSourceACLModel{{
		Read: m["read"],
	}}
}

func flattenSecretDataSourceACLOperation(ctx context.Context, acl acls.ACLDetails, respDiags *diag.Diagnostics) []SecretDataSourceACLOperationModel {
	users, diags := types.SetValueFrom(ctx, types.StringType, acl.Users)
	respDiags.Append(diags...)
	return []SecretDataSourceACLOperationModel{
		{
			CreatedAt:     types.StringValue(acl.Created.UTC().Format(time.RFC3339)),
			ProjectAccess: types.BoolValue(acl.ProjectAccess),
			UpdatedAt:     types.StringValue(acl.Updated.UTC().Format(time.RFC3339)),
			Users:         users,
		},
	}
}

func expandDateFilter(date string) *secrets.DateQuery {
	// error checks are not necessary, since they were validated by terraform validate functions
	var parts []string
	if regexp.MustCompile("^" + strings.Join(dateFilters(), "|") + ":").Match([]byte(date)) {
		parts = strings.SplitN(date, ":", 2)
	} else {
		parts = []string{date}
	}

	var parsedTime time.Time
	var filter *secrets.DateQuery

	if len(parts) == 2 {
		parsedTime, _ = time.Parse(time.RFC3339, parts[1])
		filter = &secrets.DateQuery{Date: parsedTime, Filter: secrets.DateFilter(parts[0])}
	} else {
		parsedTime, _ = time.Parse(time.RFC3339, parts[0])
		filter = &secrets.DateQuery{Date: parsedTime}
	}

	if parsedTime == (time.Time{}) {
		return nil
	}

	return filter
}
