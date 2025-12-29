package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/backup"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/cdn"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/dataplatform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/db"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/dc"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/iam"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/images"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	wrapper "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/providerwrapper/framework"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/keymanager"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/mlplatform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/monitoring"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/networking"
)

var (
	_ provider.Provider = (*vkcsProvider)(nil)
)

// Provider is a helper function to simplify provider server and testing implementation.
func Provider() provider.Provider {
	return wrapper.NewProviderWrapper(ProviderBase())
}

func ProviderBase() provider.Provider {
	return &vkcsProvider{}
}

// vkcsProvider is the provider implementation.
type vkcsProvider struct{}

type vkcsProviderModel struct {
	AuthURL                   types.String `tfsdk:"auth_url"`
	AccessToken               types.String `tfsdk:"access_token"`
	ProjectID                 types.String `tfsdk:"project_id"`
	Password                  types.String `tfsdk:"password"`
	Username                  types.String `tfsdk:"username"`
	UserDomainID              types.String `tfsdk:"user_domain_id"`
	UserDomainName            types.String `tfsdk:"user_domain_name"`
	Region                    types.String `tfsdk:"region"`
	CloudContainersAPIVersion types.String `tfsdk:"cloud_containers_api_version"`
	EndpointOverrides         types.Set    `tfsdk:"endpoint_overrides"`
	SkipClientAuth            types.Bool   `tfsdk:"skip_client_auth"`
}

type vkcsProviderEndpointOverridesModel struct {
	Backup               types.String `tfsdk:"backup"`
	BlockStorage         types.String `tfsdk:"block_storage"`
	CDN                  types.String `tfsdk:"cdn"`
	Compute              types.String `tfsdk:"compute"`
	ContainerInfra       types.String `tfsdk:"container_infra"`
	ContainerInfraAddons types.String `tfsdk:"container_infra_addons"`
	Database             types.String `tfsdk:"database"`
	DataPlatform         types.String `tfsdk:"data_platform"`
	IAMServiceUsers      types.String `tfsdk:"iam_service_users"`
	ICS                  types.String `tfsdk:"ics"`
	Image                types.String `tfsdk:"image"`
	KeyManager           types.String `tfsdk:"key_manager"`
	LoadBalancer         types.String `tfsdk:"load_balancer"`
	MLPlatform           types.String `tfsdk:"ml_platform"`
	Networking           types.String `tfsdk:"networking"`
	PublicDNS            types.String `tfsdk:"public_dns"`
	SharedFilesystem     types.String `tfsdk:"shared_filesystem"`
	Templater            types.String `tfsdk:"templater"`
}

// Metadata returns the provider type name.
func (p *vkcsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "vkcs"
}

// Schema defines the provider-level schema for configuration data.
func (p *vkcsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"auth_url": schema.StringAttribute{
				Optional:    true,
				Description: "The Identity authentication URL.",
			},
			"access_token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "A temporary token to use for authentication. You alternatively can use `OS_AUTH_TOKEN` environment variable. If both are specified, this attribute takes precedence. _note_ The token will not be renewed and will eventually expire, usually after 1 hour. If access is needed for longer than a token's lifetime, use credentials-based authentication.",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "The ID of Project to login with.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Password to login with.",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "User name to login with.",
			},
			"user_domain_id": schema.StringAttribute{
				Optional:    true,
				Description: "The id of the domain where the user resides.",
			},
			"user_domain_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the domain where the user resides.",
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Description: "A region to use.",
			},
			"cloud_containers_api_version": schema.StringAttribute{
				Optional:    true,
				Description: "Cloud Containers API version to use. _note_ Only for custom VKCS deployments.",
			},
			"skip_client_auth": schema.BoolAttribute{
				Optional:    true,
				Description: "Skip authentication on client initialization. Only applicablie if `access_token` is provided. _note_ If set to true, the endpoint catalog will not be used for discovery and all required endpoints must be provided via `endpoint_overrides`.",
				Validators: []validator.Bool{
					boolvalidator.AlsoRequires(path.MatchRoot("access_token")),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"endpoint_overrides": schema.SetNestedBlock{
				Description: "Custom endpoints for corresponding APIs. If not specified, endpoints provided by the catalog will be used.",
				Validators: []validator.Set{
					setvalidator.SizeBetween(1, 1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"backup": schema.StringAttribute{
							Optional:    true,
							Description: "Backup API custom endpoint.",
						},
						"block_storage": schema.StringAttribute{
							Optional:    true,
							Description: "Block Storage API custom endpoint.",
						},
						"cdn": schema.StringAttribute{
							Optional:    true,
							Description: "CDN API custom endpoint.",
						},
						"compute": schema.StringAttribute{
							Optional:    true,
							Description: "Compute API custom endpoint.",
						},
						"container_infra": schema.StringAttribute{
							Optional:    true,
							Description: "Cloud Containers API custom endpoint.",
						},
						"container_infra_addons": schema.StringAttribute{
							Optional:    true,
							Description: "Cloud Containers Addons API custom endpoint.",
						},
						"database": schema.StringAttribute{
							Optional:    true,
							Description: "Database API custom endpoint.",
						},
						"data_platform": schema.StringAttribute{
							Optional:    true,
							Description: "Data Platform API custom endpoint.",
						},
						"iam_service_users": schema.StringAttribute{
							Optional:    true,
							Description: "IAM Service Users API custom endpoint.",
						},
						"ics": schema.StringAttribute{
							Optional:    true,
							Description: "ICS API custom endpoint.",
						},
						"image": schema.StringAttribute{
							Optional:    true,
							Description: "Image API custom endpoint.",
						},
						"key_manager": schema.StringAttribute{
							Optional:    true,
							Description: "Key Manager API custom endpoint.",
						},
						"load_balancer": schema.StringAttribute{
							Optional:    true,
							Description: "Load Balancer API custom endpoint.",
						},
						"ml_platform": schema.StringAttribute{
							Optional:    true,
							Description: "ML Platform API custom endpoint.",
						},
						"networking": schema.StringAttribute{
							Optional:    true,
							Description: "Networking API custom endpoint.",
						},
						"public_dns": schema.StringAttribute{
							Optional:    true,
							Description: "Public DNS API custom endpoint.",
						},
						"shared_filesystem": schema.StringAttribute{
							Optional:    true,
							Description: "Shared Filesystem API custom endpoint.",
						},
						"templater": schema.StringAttribute{
							Optional:    true,
							Description: "Templater API custom endpoint.",
						},
					},
				},
			},
		},
	}
}

// Configure prepares VKCS API clients for data sources and resources.
func (p *vkcsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data vkcsProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var endpointOverrides map[string]any
	if !data.EndpointOverrides.IsNull() && !data.EndpointOverrides.IsUnknown() {
		eoElements := make([]types.Object, 0, len(data.EndpointOverrides.Elements()))
		resp.Diagnostics.Append(data.EndpointOverrides.ElementsAs(ctx, &eoElements, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		var eo vkcsProviderEndpointOverridesModel
		resp.Diagnostics.Append(eoElements[0].As(ctx, &eo, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		endpointOverrides = map[string]any{
			"backup":                 eo.Backup.ValueString(),
			"block-storage":          eo.BlockStorage.ValueString(),
			"cdn":                    eo.CDN.ValueString(),
			"compute":                eo.Compute.ValueString(),
			"container-infra":        eo.ContainerInfra.ValueString(),
			"container-infra-addons": eo.ContainerInfraAddons.ValueString(),
			"database":               eo.Database.ValueString(),
			"data-platform":          eo.DataPlatform.ValueString(),
			"iam-service-users":      eo.IAMServiceUsers.ValueString(),
			"ics":                    eo.ICS.ValueString(),
			"image":                  eo.Image.ValueString(),
			"key-manager":            eo.KeyManager.ValueString(),
			"load-balancer":          eo.LoadBalancer.ValueString(),
			"mlplatform":             eo.MLPlatform.ValueString(),
			"network":                eo.Networking.ValueString(),
			"public-dns":             eo.PublicDNS.ValueString(),
			"shared-filesystem":      eo.SharedFilesystem.ValueString(),
			"templater":              eo.Templater.ValueString(),
		}
	}

	opts := clients.ConfigOpts{
		Token:                        data.AccessToken.ValueString(),
		Username:                     data.Username.ValueString(),
		Password:                     data.Password.ValueString(),
		ProjectID:                    data.ProjectID.ValueString(),
		UserDomainID:                 data.UserDomainID.ValueString(),
		UserDomainName:               data.UserDomainName.ValueString(),
		Region:                       data.Region.ValueString(),
		IdentityEndpoint:             data.AuthURL.ValueString(),
		EndpointOverrides:            endpointOverrides,
		ContainerInfraV1MicroVersion: data.CloudContainersAPIVersion.ValueString(),
		SkipAuth:                     data.SkipClientAuth.ValueBool(),
	}

	config, err := opts.LoadAndValidate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to validate VKCS provider configuration",
			err.Error(),
		)
		return
	}

	resp.DataSourceData = config
	resp.ResourceData = config
}

// DataSources defines the data sources implemented in the provider.
func (p *vkcsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		backup.NewPlanDataSource,
		backup.NewProviderDataSource,
		backup.NewProvidersDataSource,
		cdn.NewOriginGroupDataSource,
		cdn.NewShieldingPopDataSource,
		cdn.NewShieldingPopsDataSource,
		cdn.NewSslCertificateDataSource,
		dataplatform.NewProductsDataSource,
		dataplatform.NewProductDataSource,
		dataplatform.NewTemplateDataSource,
		db.NewBackupDataSource,
		db.NewConfigGroupDataSource,
		db.NewDatastoreDataSource,
		db.NewDatastoresDataSource,
		db.NewDatastoreCapabilitiesDataSource,
		db.NewDatastoreParametersDataSource,
		dc.NewAPIOptionsDataSource,
		iam.NewServiceUserDataSource,
		iam.NewS3AccountDataSource,
		images.NewImagesDataSource,
		keymanager.NewContainerDataSource,
		keymanager.NewSecretDataSource,
		kubernetes.NewAddonDatasource,
		kubernetes.NewAddonsDatasource,
		kubernetes.NewClusterTemplatesDataSource,
		kubernetes.NewNodeGroupDataSource,
		kubernetes.NewSecurityPolicyTemplatesDataSource,
		kubernetes.NewSecurityPolicyTemplateDataSource,
		networking.NewPortDataSource,
		networking.NewSubnetDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *vkcsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		backup.NewPlanResource,
		cdn.NewOriginGroupResource,
		cdn.NewResourceResource,
		cdn.NewSslCertificateResource,
		dataplatform.NewClusterResource,
		db.NewBackupResource,
		dc.NewRouterResource,
		dc.NewInterfaceResource,
		dc.NewBGPInstanceResource,
		dc.NewBGPNeighborResource,
		dc.NewBGPStaticAnnounceResource,
		dc.NewStaticRouteResource,
		dc.NewVRRPResource,
		dc.NewVRRPInterfaceResource,
		dc.NewVRRPAddressResource,
		dc.NewConntrackHelperResource,
		dc.NewIPPortForwardingResource,
		iam.NewServiceUserResource,
		iam.NewS3AccountResource,
		kubernetes.NewAddonResource,
		kubernetes.NewSecurityPolicyResource,
		mlplatform.NewJupyterHubResource,
		mlplatform.NewMLFlowResource,
		mlplatform.NewMLFlowDeployResource,
		mlplatform.NewSparkK8SResource,
		mlplatform.NewK8SRegistryResource,
		monitoring.NewResource,
		networking.NewAnycastIPResource,
	}
}
