package vkcs

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceKubernetesClusterTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubernetesClusterTemplateRead,
		Schema: map[string]*schema.Schema{
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_template_uuid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"apiserver_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cluster_distro": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_nameserver": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"docker_storage_driver": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"docker_volume_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"external_network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_flavor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"floating_ip_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"image": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"insecure_registry": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"keypair_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"master_lb_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"network_driver": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"no_proxy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"registry_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"server_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tls_disabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"volume_driver": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deprecated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceKubernetesClusterTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	containerInfraClient, err := config.ContainerInfraV1Client(config.GetRegion())
	if err != nil {
		return diag.Errorf("error creating VKCS container infra client: %s", err)
	}
	templateIdentifierKey, err := ensureOnlyOnePresented(d, "name", "version", "cluster_template_uuid")
	if err != nil {
		return diag.FromErr(err)
	}
	templateIdentifier := d.Get(templateIdentifierKey).(string)
	var ct *clusterTemplate
	ct, err = clusterTemplateGet(containerInfraClient, templateIdentifier).Extract()
	if err != nil {
		return diag.Errorf("error getting vkcs_kubernetes_clustertemplate %s: %s", templateIdentifier, err)
	}

	d.SetId(ct.UUID)

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
	d.Set("region", getRegion(d, config))
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
