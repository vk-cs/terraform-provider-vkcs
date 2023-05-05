package vkcs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	cg "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/config_groups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/datastores"
)

func resourceDatabaseConfigGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseConfigGroupCreate,
		ReadContext:   resourceDatabaseConfigGroupRead,
		UpdateContext: resourceDatabaseConfigGroupUpdate,
		DeleteContext: resourceDatabaseConfigGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"datastore": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Version of the datastore.",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Type of the datastore.",
						},
					},
				},
				Description: "Object that represents datastore of the config group. Changing this creates a new config group.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "The name of the config group.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "The description of the config group.",
			},
			"values": {
				Type:        schema.TypeMap,
				Required:    true,
				ForceNew:    false,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Map of configuration parameters in format \"key\": \"value\".",
			},
			"updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of config group's last update",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of config group's creation",
			},
		},
		Description: "Provides a db config group resource. This can be used to create, update and delete db config group.\n" +
			"**New since v.0.1.7**.",
	}
}

func resourceDatabaseConfigGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	createOpts := cg.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	v := d.Get("datastore")
	datastore, err := extractDatabaseDatastore(v.([]interface{}))
	if err != nil {
		return diag.Errorf("unable to determine vkcs_db_config_group datastore")
	}
	createOpts.Datastore = &datastore

	v = d.Get("values")
	values, err := retrieveDatabaseConfigGroupValues(DatabaseV1Client, datastore, v.(map[string]interface{}))
	if err != nil {
		return diag.Errorf("unable to determine vkcs_db_config_group values: %s", err)
	}
	createOpts.Values = values

	log.Printf("[DEBUG] vkcs_db_backup create options: %#v", createOpts)

	configGrp := cg.ConfigGroup{
		Configuration: &createOpts,
	}

	configGroup, err := cg.Create(DatabaseV1Client, &configGrp).Extract()
	if err != nil {
		return diag.Errorf("error creating vkcs_db_config_group: %s", err)
	}

	// Store the ID now
	d.SetId(configGroup.ID)

	return resourceDatabaseConfigGroupRead(ctx, d, meta)
}

func resourceDatabaseConfigGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	configGroup, err := cg.Get(DatabaseV1Client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_db_config_group"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_db_config_group %s: %#v", d.Id(), configGroup)

	d.Set("name", configGroup.Name)
	ds := datastores.DatastoreShort{
		Type:    configGroup.DatastoreName,
		Version: configGroup.DatastoreVersionName,
	}
	d.Set("datastore", flattenDatabaseInstanceDatastore(ds))
	d.Set("values", flattenDatabaseConfigGroupValues(configGroup.Values))

	d.Set("updated", configGroup.Updated)
	d.Set("created", configGroup.Created)
	d.Set("description", configGroup.Description)

	return nil
}

func resourceDatabaseConfigGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	v := d.Get("datastore")
	datastore, err := extractDatabaseDatastore(v.([]interface{}))
	if err != nil {
		return diag.Errorf("unable to determine vkcs_db_config_group datastore")
	}

	v = d.Get("values")
	values, err := retrieveDatabaseConfigGroupValues(DatabaseV1Client, datastore, v.(map[string]interface{}))
	if err != nil {
		return diag.Errorf("unable to determine vkcs_db_config_group values: %s", err)
	}

	updateOpts := cg.UpdateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Values:      values,
	}

	log.Printf("[DEBUG] vkcs_db_config_group update options: %#v", updateOpts)
	update := cg.UpdateOpt{
		Configuration: &updateOpts,
	}

	err = cg.Update(DatabaseV1Client, d.Id(), &update).ExtractErr()
	if err != nil {
		return diag.Errorf("error updating vkcs_db_config_group: %s", err)
	}
	return resourceDatabaseConfigGroupRead(ctx, d, meta)
}

func resourceDatabaseConfigGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	err = cg.Delete(DatabaseV1Client, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error deleting vkcs_db_config_group"))
	}

	return nil
}
