package vkcs

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/acls"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/containers"
)

func dataSourceKeyManagerContainer() *schema.Resource {
	ret := &schema.Resource{
		ReadContext: dataSourceKeyManagerContainerRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
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
			},

			"container_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"creator_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"consumers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"acl": {
				Type:     schema.TypeList,
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

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	elem := &schema.Resource{
		Schema: make(map[string]*schema.Schema),
	}
	for _, aclOp := range getSupportedACLOperations() {
		elem.Schema[aclOp] = getACLSchema()
	}
	ret.Schema["acl"].Elem = elem

	return ret
}

func dataSourceKeyManagerContainerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	kmClient, err := config.KeyManagerV1Client(getRegion(d, config))
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
	d.Set("region", getRegion(d, config))

	return nil
}
