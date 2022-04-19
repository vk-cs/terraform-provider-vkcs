package vkcs

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDatabaseDatabase() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatabaseDatabaseRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"dbms_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"charset": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"collate": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func dataSourceDatabaseDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating vkcs database client: %s", err)
	}

	id := d.Get("id").(string)
	databaseID := strings.SplitN(id, "/", 2)
	if len(databaseID) != 2 {
		return diag.Errorf("invalid vkcs_db_database id: %s", id)
	}

	dbmsID := databaseID[0]
	databaseName := databaseID[1]
	dbmsResp, err := getDBMSResource(DatabaseV1Client, dbmsID)
	if err != nil {
		return diag.Errorf("error while getting resource: %s", err)
	}
	var dbmsType string
	if _, ok := dbmsResp.(instanceResp); ok {
		dbmsType = dbmsTypeInstance
	}
	if _, ok := dbmsResp.(dbClusterResp); ok {
		dbmsType = dbmsTypeCluster
	}
	exists, err := databaseDatabaseExists(DatabaseV1Client, dbmsID, databaseName, dbmsType)
	if err != nil {
		return diag.Errorf("error checking if vkcs_db_database %s exists: %s", d.Id(), err)
	}

	if !exists {
		d.SetId("")
		return nil
	}

	d.SetId(id)
	d.Set("name", databaseName)
	d.Set("dbms_id", dbmsID)
	return nil
}
