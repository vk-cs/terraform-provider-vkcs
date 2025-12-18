package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &vkcsProvider{}
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
		},
	}
}

// Configure prepares VKCS API clients for data sources and resources.
func (p *vkcsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	config, diags := clients.ConfigureProvider(ctx, req)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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
