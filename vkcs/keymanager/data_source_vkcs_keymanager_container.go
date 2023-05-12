package keymanager

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/acls"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/containers"
)

func DataSourceKeyManagerContainer() *schema.Resource {
	ret := &schema.Resource{
		ReadContext: dataSourceKeyManagerContainerRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region in which to obtain the KeyManager client. A KeyManager client is needed to fetch a container. If omitted, the `region` argument of the provider is used.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Container name.",
			},

			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The container type.",
			},

			"secret_refs": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"secret_ref": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Description: "A set of dictionaries containing references to secrets.",
			},

			"container_ref": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The container reference / where to find the container.",
			},

			"creator_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creator of the container.",
			},

			"consumers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the consumer.",
						},
						"url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The consumer URL.",
						},
					},
				},
				Description: "The list of the container consumers.",
			},

			"acl": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of ACLs assigned to a container.",
			},

			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the container was created.",
			},

			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the container was last updated.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the container.",
			},
		},
		Description: "Use this data source to get the ID of an available Key container.",
	}

	elem := &schema.Resource{
		Schema: make(map[string]*schema.Schema),
	}
	for _, aclOp := range getSupportedACLOperations() {
		elem.Schema[aclOp] = getACLSchema()
		elem.Schema[aclOp].Description = fmt.Sprintf("Block that describes %s operation.", aclOp)
	}
	ret.Schema["acl"].Elem = elem

	return ret
}

func dataSourceKeyManagerContainerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	kmClient, err := config.KeyManagerV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS keymanager client: %s", err)
	}

	listOpts := containers.ListOpts{
		Name: d.Get("name").(string),
	}

	log.Printf("[DEBUG] Containers List Options: %#v", listOpts)

	allPages, err := containers.List(kmClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to query vkcs_keymanager_container containers: %s", err)
	}

	allContainers, err := containers.ExtractContainers(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve vkcs_keymanager_container containers: %s", err)
	}

	if len(allContainers) < 1 {
		return diag.Errorf("Your query returned no vkcs_keymanager_container results. " +
			"Please change your search criteria and try again.")
	}

	if len(allContainers) > 1 {
		log.Printf("[DEBUG] Multiple vkcs_keymanager_container results found: %#v", allContainers)
		return diag.Errorf("Your query returned more than one result. Please try a more " +
			"specific search criteria.")
	}

	container := allContainers[0]

	log.Printf("[DEBUG] Retrieved vkcs_keymanager_container %s: %#v", d.Id(), container)

	uuid := keyManagerContainerGetUUIDfromContainerRef(container.ContainerRef)

	d.SetId(uuid)
	d.Set("name", container.Name)

	d.Set("creator_id", container.CreatorID)
	d.Set("container_ref", container.ContainerRef)
	d.Set("type", container.Type)
	d.Set("status", container.Status)
	d.Set("created_at", container.Created.Format(time.RFC3339))
	d.Set("updated_at", container.Updated.Format(time.RFC3339))
	d.Set("consumers", flattenKeyManagerContainerConsumers(container.Consumers))

	d.Set("secret_refs", flattenKeyManagerContainerSecretRefs(container.SecretRefs))

	acl, err := acls.GetContainerACL(kmClient, d.Id()).Extract()
	if err != nil {
		log.Printf("[DEBUG] Unable to get %s container acls: %s", uuid, err)
	}
	d.Set("acl", flattenKeyManagerACLs(acl))

	// Set the region
	d.Set("region", util.GetRegion(d, config))

	return nil
}
