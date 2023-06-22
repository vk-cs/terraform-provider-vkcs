package keymanager

import (
	"context"
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/acls"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/containers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
)

var (
	_ datasource.DataSource              = &ContainerDataSource{}
	_ datasource.DataSourceWithConfigure = &ContainerDataSource{}
)

func NewContainerDataSource() datasource.DataSource {
	return &ContainerDataSource{}
}

type ContainerDataSource struct {
	config clients.Config
}

type ContainerDataSourceModel struct {
	ID     types.String `tfsdk:"id"`
	Region types.String `tfsdk:"region"`

	ACL          []ContainerDataSourceACLModel       `tfsdk:"acl"`
	Consumers    []ContainerDataSourceConsumerModel  `tfsdk:"consumers"`
	ContainerRef types.String                        `tfsdk:"container_ref"`
	CreatedAt    types.String                        `tfsdk:"created_at"`
	CreatorID    types.String                        `tfsdk:"creator_id"`
	Name         types.String                        `tfsdk:"name"`
	SecretRefs   []ContainerDataSourceSecretRefModel `tfsdk:"secret_refs"`
	Status       types.String                        `tfsdk:"status"`
	Type         types.String                        `tfsdk:"type"`
	UpdatedAt    types.String                        `tfsdk:"updated_at"`
}

type ContainerDataSourceACLModel struct {
	Read []ContainerDataSourceACLOperationModel `tfsdk:"read"`
}

type ContainerDataSourceACLOperationModel struct {
	CreatedAt     types.String `tfsdk:"created_at"`
	ProjectAccess types.Bool   `tfsdk:"project_access"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
	Users         types.Set    `tfsdk:"users"`
}

type ContainerDataSourceConsumerModel struct {
	Name types.String `tfsdk:"name"`
	URL  types.String `tfsdk:"url"`
}

type ContainerDataSourceSecretRefModel struct {
	Name      types.String `tfsdk:"name"`
	SecretRef types.String `tfsdk:"secret_ref"`
}

func (d *ContainerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_keymanager_container"
}

func (d *ContainerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource.",
			},

			"region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The region in which to obtain the VKCS Key Manager client. If omitted, the `region` argument of the provider is used.",
			},

			"acl": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"read": schema.ListNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"created_at": schema.StringAttribute{
										Computed:    true,
										Description: "The date the container ACL was created.",
									},

									"project_access": schema.BoolAttribute{
										Computed:    true,
										Description: "Whether the container is accessible project wide.",
									},

									"updated_at": schema.StringAttribute{
										Computed:    true,
										Description: "The date the container ACL was last updated.",
									},

									"users": schema.SetAttribute{
										ElementType: types.StringType,
										Computed:    true,
										Description: "The list of user IDs, which are allowed to access the container, when `project_access` is set to `false`.",
									},
								},
							},
							Computed:    true,
							Description: "Object that describes read operation.",
						},
					},
				},
				Computed:    true,
				Description: "ACLs assigned to a container.",
			},

			"consumers": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Optional:    true,
							Description: "The name of the consumer.",
						},

						"url": schema.StringAttribute{
							Optional:    true,
							Description: "The consumer URL.",
						},
					},
				},
				Computed:    true,
				Description: "The list of the container consumers.",
			},

			"container_ref": schema.StringAttribute{
				Computed:    true,
				Description: "The container reference / where to find the container.",
			},

			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "The date the container was created.",
			},

			"creator_id": schema.StringAttribute{
				Computed:    true,
				Description: "The creator of the container.",
			},

			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The Container name.",
			},

			"secret_refs": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Optional: true,
						},

						"secret_ref": schema.StringAttribute{
							Optional: true,
						},
					},
				},
				Computed:    true,
				Description: "A set of dictionaries containing references to secrets.",
			},

			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the container.",
			},

			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The container type.",
			},

			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "The date the container was last updated.",
			},
		},
		Description: "Use this data source to get the ID of an available Key container.",
	}
}

func (d *ContainerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *ContainerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ContainerDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}

	client, err := d.config.KeyManagerV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS Key Manager API client", err.Error())
		return
	}

	listOpts := containers.ListOpts{
		Name: data.Name.ValueString(),
	}

	tflog.Debug(ctx, "Calling Key Manager API to list containers", map[string]interface{}{"list_opts": fmt.Sprintf("%#v", listOpts)})

	allPages, err := containers.List(client, &listOpts).AllPages()
	if err != nil {
		resp.Diagnostics.AddError("Error calling VKCS Key Manager API", err.Error())
		return
	}

	allContainers, err := containers.ExtractContainers(allPages)
	if err != nil {
		resp.Diagnostics.AddError("Error reading VKCS Key Manager API response", err.Error())
		return
	}

	tflog.Debug(ctx, "Called Key Manager API to list containers", map[string]interface{}{"all_containers_len": len(allContainers)})

	if len(allContainers) < 1 {
		resp.Diagnostics.AddError("Your query returned no results",
			"Please change your search criteria and try again")
	}

	if len(allContainers) > 1 {
		resp.Diagnostics.AddError("Your query returned more than one result",
			"Please try a more specific search criteria")
	}

	container := allContainers[0]
	tflog.Debug(ctx, "Retrieved the container", map[string]interface{}{"container": fmt.Sprintf("%#v", container)})

	id := GetUUIDfromContainerRef(container.ContainerRef)
	ctx = tflog.SetField(ctx, "id", id)

	data.ID = types.StringValue(id)
	data.Region = types.StringValue(region)
	data.Consumers = flattenContainerDataSourceConsumers(ctx, container.Consumers, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ContainerRef = types.StringValue(container.ContainerRef)
	data.CreatedAt = types.StringValue(container.Created.Format(time.RFC3339))
	data.CreatorID = types.StringValue(container.CreatorID)
	data.Name = types.StringValue(container.Name)
	data.SecretRefs = flattenContainerDataSourceSecretRefs(ctx, container.SecretRefs, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Status = types.StringValue(container.Status)
	data.Type = types.StringValue(container.Type)
	data.UpdatedAt = types.StringValue(container.Updated.Format(time.RFC3339))

	tflog.Debug(ctx, "Calling Key Manager API to get container's acls")

	acl, err := acls.GetContainerACL(client, id).Extract()
	if err != nil {
		tflog.Debug(ctx, "Error calling Key Manager API to get container's acls", map[string]interface{}{"error": err.Error()})
	} else {
		tflog.Debug(ctx, "Called Key Manager API to get container's acls", map[string]interface{}{"acls": fmt.Sprintf("%#v", acl)})
	}

	data.ACL = flattenContainerDataSourceACL(ctx, acl, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenContainerDataSourceACL(ctx context.Context, in *acls.ACL, respDiags *diag.Diagnostics) []ContainerDataSourceACLModel {
	r := []ContainerDataSourceACLModel{}

	if in == nil {
		return r
	}

	acl := ContainerDataSourceACLModel{}
	if v, ok := (*in)["read"]; ok {
		acl.Read = flattenContainerDataSourceACLOperation(ctx, v, respDiags)
	}

	if len(acl.Read) > 0 {
		r = append(r, acl)
	}

	return r
}

func flattenContainerDataSourceACLOperation(ctx context.Context, in acls.ACLDetails, respDiags *diag.Diagnostics) []ContainerDataSourceACLOperationModel {
	users, diags := types.SetValueFrom(ctx, types.StringType, in.Users)
	respDiags.Append(diags...)

	return []ContainerDataSourceACLOperationModel{
		{
			CreatedAt:     types.StringValue(in.Created.UTC().Format(time.RFC3339)),
			ProjectAccess: types.BoolValue(in.ProjectAccess),
			UpdatedAt:     types.StringValue(in.Updated.UTC().Format(time.RFC3339)),
			Users:         users,
		},
	}
}

func flattenContainerDataSourceConsumers(_ context.Context, in []containers.ConsumerRef, _ *diag.Diagnostics) []ContainerDataSourceConsumerModel {
	r := make([]ContainerDataSourceConsumerModel, len(in))
	for i, cr := range in {
		r[i] = ContainerDataSourceConsumerModel{
			Name: types.StringValue(cr.Name),
			URL:  types.StringValue(cr.URL),
		}
	}
	return r
}

func flattenContainerDataSourceSecretRefs(_ context.Context, in []containers.SecretRef, _ *diag.Diagnostics) []ContainerDataSourceSecretRefModel {
	r := make([]ContainerDataSourceSecretRefModel, len(in))
	for i, sr := range in {
		r[i] = ContainerDataSourceSecretRefModel{
			Name:      types.StringValue(sr.Name),
			SecretRef: types.StringValue(sr.SecretRef),
		}
	}
	return r
}
