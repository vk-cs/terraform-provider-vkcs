package kubernetes

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfra/v1/clustertemplates"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func DataSourceKubernetesClusterTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubernetesClusterTemplateRead,
		Schema: map[string]*schema.Schema{
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Kubernetes version of the cluster. _note_ Only one of `name`, `version` or `id` must be specified.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The UUID of the cluster template. _note_ Only one of `name`, `version` or `id` must be specified.",
			},
			"cluster_template_uuid": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
				Deprecated:    "This argument is deprecated, please, use the `id` attribute instead.",
				Description:   "The UUID of the cluster template. _note_ Only one of `name`, `version`, or `cluster_template_uuid` must be specified.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region in which to obtain the V1 Container Infra client. If omitted, the `region` argument of the provider is used.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the cluster template. _note_ Only one of `name`, `version` or `id` must be specified.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The project of the cluster template.",
			},
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user of the cluster template.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time at which cluster template was created.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time at which cluster template was updated.",
			},
			"apiserver_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The API server port for the Container Orchestration Engine for this cluster template.",
			},
			"cluster_distro": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The distro for the cluster (fedora-atomic, coreos, etc.).",
			},
			"dns_nameserver": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Address of the DNS nameserver that is used in nodes of the cluster.",
			},
			"docker_storage_driver": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Docker storage driver. Changing this updates the Docker storage driver of the existing cluster template.",
			},
			"docker_volume_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size (in GB) of the Docker volume.",
			},
			"external_network_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the external network that will be used for the cluster.",
			},
			"flavor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of flavor for the nodes of the cluster.",
			},
			"master_flavor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of flavor for the master nodes.",
			},
			"floating_ip_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether created cluster should create IP floating IP for every node or not.",
			},
			"image": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The reference to an image that is used for nodes of the cluster.",
			},
			"insecure_registry": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The insecure registry URL for the cluster template.",
			},
			"keypair_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the Compute service SSH keypair.",
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The list of key value pairs representing additional properties of the cluster template.",
			},
			"master_lb_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether created cluster should has a loadbalancer for master nodes or not.",
			},
			"network_driver": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the driver for the container network.",
			},
			"no_proxy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A comma-separated list of IP addresses that shouldn't be used in the cluster.",
			},
			"public": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether cluster template should be public.",
			},
			"registry_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether Docker registry is enabled in the cluster.",
			},
			"server_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The server type for the cluster template.",
			},
			"tls_disabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the TLS should be disabled in the cluster.",
			},
			"volume_driver": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the driver that is used for the volumes of the cluster nodes.",
			},
			"deprecated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time at which the cluster template is deprecated.",
			},
		},
		Description: "Use this data source to get the ID of an available VKCS kubernetes cluster template.",
	}
}

func dataSourceKubernetesClusterTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	containerInfraClient, err := config.ContainerInfraV1Client(config.GetRegion())
	if err != nil {
		return diag.Errorf("error creating VKCS container infra client: %s", err)
	}
	templateIdentifierKey, err := util.EnsureOnlyOnePresented(d, "name", "version", "id", "cluster_template_uuid")
	if err != nil {
		return diag.FromErr(err)
	}
	templateIdentifier := d.Get(templateIdentifierKey).(string)
	var ct *clustertemplates.ClusterTemplate
	ct, err = clustertemplates.Get(containerInfraClient, templateIdentifier).Extract()
	if err != nil {
		return diag.Errorf("error getting vkcs_kubernetes_clustertemplate %s: %s", templateIdentifier, err)
	}

	d.SetId(ct.UUID)
	d.Set("cluster_template_uuid", ct.UUID)
	d.Set("project_id", ct.ProjectID)
	d.Set("user_id", ct.UserID)
	d.Set("apiserver_port", ct.APIServerPort)
	d.Set("cluster_distro", ct.ClusterDistro)
	d.Set("dns_nameserver", ct.DNSNameServer)
	d.Set("docker_storage_driver", ct.DockerStorageDriver)
	d.Set("docker_volume_size", ct.DockerVolumeSize)
	d.Set("external_network_id", ct.ExternalNetworkID)
	d.Set("flavor", ct.FlavorID)
	d.Set("master_flavor", ct.MasterFlavorID)
	d.Set("floating_ip_enabled", ct.FloatingIPEnabled)
	d.Set("image", ct.ImageID)
	d.Set("insecure_registry", ct.InsecureRegistry)
	d.Set("keypair_id", ct.KeyPairID)
	d.Set("labels", ct.Labels)
	d.Set("master_lb_enabled", ct.MasterLBEnabled)
	d.Set("network_driver", ct.NetworkDriver)
	d.Set("no_proxy", ct.NoProxy)
	d.Set("public", ct.Public)
	d.Set("registry_enabled", ct.RegistryEnabled)
	d.Set("server_type", ct.ServerType)
	d.Set("tls_disabled", ct.TLSDisabled)
	d.Set("volume_driver", ct.VolumeDriver)
	d.Set("name", ct.Name)
	d.Set("version", ct.Version)
	d.Set("region", util.GetRegion(d, config))
	d.Set("deprecated_at", "")

	if err := d.Set("created_at", ct.CreatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set vkcs_containerinfra_clustertemplate created_at: %s", err)
	}

	if err := d.Set("updated_at", ct.UpdatedAt.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set vkcs_containerinfra_clustertemplate updated_at: %s", err)
	}

	if !ct.DeprecatedAt.IsZero() {
		if err := d.Set("deprecated_at", ct.DeprecatedAt.Format(time.RFC3339)); err != nil {
			log.Printf("[DEBUG] Unable to set vkcs_containerinfra_clustertemplate deprecated_at: %s", err)
		}
	}

	return nil
}
