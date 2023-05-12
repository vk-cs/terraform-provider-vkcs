package db

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/instances"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func DataSourceDatabaseUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatabaseUserRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the user in form \"dbms_id/user_name\".",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the user. Changing this creates a new user.",
			},

			"dbms_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the instance or cluster that user is created for.",
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The password of the user.",
			},

			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP address of the host that user will be accessible from.",
			},

			"databases": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of names of the databases, that user is created for.",
			},
		},
		Description: "Use this data source to get the information on a db user resource.",
	}
}

func dataSourceDatabaseUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	DatabaseV1Client, err := config.DatabaseV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating vkcs database client: %s", err)
	}

	id := d.Get("id").(string)
	userID := strings.SplitN(id, "/", 2)
	if len(userID) != 2 {
		return diag.Errorf("invalid vkcs_db_user id: %s", id)
	}

	dbmsID := userID[0]
	userName := userID[1]
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
	exists, userObj, err := databaseUserExists(DatabaseV1Client, dbmsID, userName, dbmsType)
	if err != nil {
		return diag.Errorf("error checking if vkcs_db_user %s exists: %s", d.Id(), err)
	}

	if !exists {
		d.SetId("")
		return nil
	}

	d.SetId(id)
	d.Set("name", userName)

	databases := flattenDatabaseUserDatabases(userObj.Databases)
	if err := d.Set("databases", databases); err != nil {
		return diag.Errorf("unable to set databases: %s", err)
	}
	d.Set("dbms_id", dbmsID)
	return nil
}
