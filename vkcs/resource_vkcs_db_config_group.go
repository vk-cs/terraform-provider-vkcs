package vkcs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"values": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDatabaseConfigGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	createOpts := dbConfigGroupCreateOpts{
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

	configGrp := dbConfigGroup{
		Configuration: &createOpts,
	}

	configGroup, err := dbConfigGroupCreate(DatabaseV1Client, &configGrp).extract()
	if err != nil {
		return diag.Errorf("error creating vkcs_db_config_group: %s", err)
	}

	// Store the ID now
	d.SetId(configGroup.ID)

	return resourceDatabaseConfigGroupRead(ctx, d, meta)
}

func resourceDatabaseConfigGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	configGroup, err := dbConfigGroupGet(DatabaseV1Client, d.Id()).extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_db_config_group"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_db_config_group %s: %#v", d.Id(), configGroup)

	d.Set("name", configGroup.Name)
	ds := dataStore{
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
	config := meta.(configer)
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

	updateOpts := dbConfigGroupUpdateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Values:      values,
	}

	log.Printf("[DEBUG] vkcs_db_config_group update options: %#v", updateOpts)
	update := dbConfigGroupUpdateOpt{
		Configuration: &updateOpts,
	}

	err = dbConfigGroupUpdate(DatabaseV1Client, d.Id(), &update).ExtractErr()
	if err != nil {
		return diag.Errorf("error updating vkcs_db_config_group: %s", err)
	}
	return resourceDatabaseConfigGroupRead(ctx, d, meta)
}

func resourceDatabaseConfigGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	err = dbConfigGroupDelete(DatabaseV1Client, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error deleting vkcs_db_config_group"))
	}

	return nil
}
