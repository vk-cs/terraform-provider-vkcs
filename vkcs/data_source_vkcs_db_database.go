package vkcs

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/instances"
)

func dataSourceDatabaseDatabase() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatabaseDatabaseRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the database in form \"dbms_id/db_name\".",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the database.",
			},

			"dbms_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the instance or cluster that database is created for.",
			},

			"charset": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Type of charset used for the database.",
			},

			"collate": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Collate option of the database.",
			},
		},
		Description: "Use this data source to get the information on a db database resource.",
	}
}

func dataSourceDatabaseDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
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
	if _, ok := dbmsResp.(instances.InstanceResp); ok {
		dbmsType = db.DBMSTypeInstance
	}
	if _, ok := dbmsResp.(clusters.ClusterResp); ok {
		dbmsType = db.DBMSTypeCluster
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
