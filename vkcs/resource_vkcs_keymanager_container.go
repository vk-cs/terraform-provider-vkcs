package vkcs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/acls"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/containers"
)

func resourceKeyManagerContainer() *schema.Resource {
	ret := &schema.Resource{
		CreateContext: resourceKeyManagerContainerCreate,
		ReadContext:   resourceKeyManagerContainerRead,
		UpdateContext: resourceKeyManagerContainerUpdate,
		DeleteContext: resourceKeyManagerContainerDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the KeyManager client. A KeyManager client is needed to create a container. If omitted, the `region` argument of the provider is used. Changing this creates a new container.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Human-readable name for the Container. Does not have to be unique.",
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"generic", "rsa", "certificate",
				}, false),
				Description: "Used to indicate the type of container. Must be one of `generic`, `rsa` or `certificate`.",
			},

			"secret_refs": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the secret reference. The reference names must correspond the container type, more details are available [here](https://docs.openstack.org/barbican/stein/api/reference/containers.html).",
						},
						"secret_ref": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The secret reference / where to find the secret, URL.",
						},
					},
				},
				Description: "A set of dictionaries containing references to secrets. The structure is described below.",
			},

			"acl": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Allows to control an access to a container. Currently only the `read` operation is supported. If not specified, the container is accessible project wide. The `read` structure is described below.",
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
				Description: "The list of the container consumers. The structure is described below.",
			},

			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the container ACL was created.",
			},

			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the container ACL was last updated.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the container.",
			},
		},
		Description: "Manages a key container resource within VKCS.",
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

func resourceKeyManagerContainerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	kmClient, err := config.KeyManagerV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS KeyManager client: %s", err)
	}

	containerType := keyManagerContainerType(d.Get("type").(string))

	createOpts := containers.CreateOpts{
		Name:       d.Get("name").(string),
		Type:       containerType,
		SecretRefs: expandKeyManagerContainerSecretRefs(d.Get("secret_refs").(*schema.Set)),
	}

	log.Printf("[DEBUG] Create Options for vkcs_keymanager_container: %#v", createOpts)

	container, err := containers.Create(kmClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating vkcs_keymanager_container: %s", err)
	}

	uuid := keyManagerContainerGetUUIDfromContainerRef(container.ContainerRef)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"ACTIVE"},
		Refresh:    keyManagerContainerWaitForContainerCreation(kmClient, uuid),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vkcs_keymanager_container: %s", err)
	}

	d.SetId(uuid)

	d.Partial(true)

	// set the acl first before setting the secret refs
	if _, ok := d.GetOk("acl"); ok {
		setOpts := expandKeyManagerACLs(d.Get("acl"))
		_, err = acls.SetContainerACL(kmClient, uuid, setOpts).Extract()
		if err != nil {
			return diag.Errorf("Error settings ACLs for the vkcs_keymanager_container: %s", err)
		}
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vkcs_keymanager_container: %s", err)
	}

	d.Partial(false)

	return resourceKeyManagerContainerRead(ctx, d, meta)
}

func resourceKeyManagerContainerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	kmClient, err := config.KeyManagerV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS keymanager client: %s", err)
	}

	container, err := containers.Get(kmClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_keymanager_container"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_keymanager_container %s: %#v", d.Id(), container)

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
		log.Printf("[DEBUG] Unable to get %s container acls: %s", d.Id(), err)
	}
	d.Set("acl", flattenKeyManagerACLs(acl))

	// Set the region
	d.Set("region", getRegion(d, config))

	return nil
}

func resourceKeyManagerContainerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	kmClient, err := config.KeyManagerV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS keymanager client: %s", err)
	}

	if d.HasChange("acl") {
		updateOpts := expandKeyManagerACLs(d.Get("acl"))
		_, err := acls.UpdateContainerACL(kmClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating vkcs_keymanager_container %s acl: %s", d.Id(), err)
		}
	}

	return resourceKeyManagerContainerRead(ctx, d, meta)
}

func resourceKeyManagerContainerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	kmClient, err := config.KeyManagerV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS keymanager client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"DELETED"},
		Refresh:    keyManagerContainerWaitForContainerDeletion(kmClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
