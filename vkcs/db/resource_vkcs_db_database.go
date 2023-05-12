package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/databases"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/instances"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func ResourceDatabaseDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseDatabaseCreate,
		ReadContext:   resourceDatabaseDatabaseRead,
		DeleteContext: resourceDatabaseDatabaseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(dbDatabaseCreateTimeout),
			Delete: schema.DefaultTimeout(dbDatabaseDeleteTimeout),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the database. Changing this creates a new database.",
			},

			"dbms_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the instance or cluster that database is created for.",
			},

			"charset": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Type of charset used for the database. Changing this creates a new database.",
			},

			"collate": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Collate option of the database.  Changing this creates a new database.",
			},

			"dbms_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of dbms for the database, can be \"instance\" or \"cluster\".",
			},
		},
		Description: "Provides a db database resource. This can be used to create, modify and delete db databases.",
	}
}

func resourceDatabaseDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	DatabaseV1Client, err := config.DatabaseV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	databaseName := d.Get("name").(string)
	dbmsID := d.Get("dbms_id").(string)

	dbmsResp, err := getDBMSResource(DatabaseV1Client, dbmsID)
	if err != nil {
		return diag.Errorf("error while getting instance or cluster: %s", err)
	}
	var dbmsType string
	if instanceResource, ok := dbmsResp.(*instances.InstanceResp); ok {
		if util.IsOperationNotSupported(instanceResource.DataStore.Type, Redis, Tarantool) {
			return diag.Errorf("operation not supported for this datastore")
		}
		if instanceResource.ReplicaOf != nil {
			return diag.Errorf("operation not supported for replica")
		}
		dbmsType = db.DBMSTypeInstance
	}
	if clusterResource, ok := dbmsResp.(*clusters.ClusterResp); ok {
		if util.IsOperationNotSupported(clusterResource.DataStore.Type, Redis, Tarantool) {
			return diag.Errorf("operation not supported for this datastore")
		}
		dbmsType = db.DBMSTypeCluster
	}
	var databasesList databases.BatchCreateOpts

	db := databases.CreateOpts{
		Name:    databaseName,
		CharSet: d.Get("charset").(string),
		Collate: d.Get("collate").(string),
	}

	databasesList.Databases = append(databasesList.Databases, db)
	err = databases.Create(DatabaseV1Client, dbmsID, &databasesList, dbmsType).ExtractErr()
	if err != nil {
		return diag.Errorf("error creating vkcs_db_database: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    databaseDatabaseStateRefreshFunc(DatabaseV1Client, dbmsID, databaseName, dbmsType),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      dbDatabaseDelay,
		MinTimeout: dbDatabaseMinTimeout,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for vkcs_db_database %s to be created: %s", databaseName, err)
	}

	// Store the ID now
	d.SetId(fmt.Sprintf("%s/%s", dbmsID, databaseName))
	// Store dbms type
	d.Set("dbms_type", dbmsType)

	return resourceDatabaseDatabaseRead(ctx, d, meta)
}

func resourceDatabaseDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	DatabaseV1Client, err := config.DatabaseV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating vkcs database client: %s", err)
	}

	databaseID := strings.SplitN(d.Id(), "/", 2)
	if len(databaseID) != 2 {
		return diag.Errorf("invalid vkcs_db_database ID: %s", d.Id())
	}

	dbmsID := databaseID[0]
	databaseName := databaseID[1]

	var dbmsType string
	if dbmsTypeRaw, ok := d.GetOk("dbms_type"); ok {
		dbmsType = dbmsTypeRaw.(string)
	} else {
		dbmsType = db.DBMSTypeInstance
	}

	_, err = getDBMSResource(DatabaseV1Client, dbmsID)
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error retrieving vkcs_db_database"))
	}

	exists, err := databaseDatabaseExists(DatabaseV1Client, dbmsID, databaseName, dbmsType)
	if err != nil {
		return diag.Errorf("error checking if vkcs_db_database %s exists: %s", d.Id(), err)
	}

	if !exists {
		d.SetId("")
		return nil
	}

	d.Set("name", databaseName)
	d.Set("dbms_id", dbmsID)
	d.Set("dbms_type", dbmsType)

	return nil
}

func resourceDatabaseDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	DatabaseV1Client, err := config.DatabaseV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating vkcs database client: %s", err)
	}

	databaseID := strings.SplitN(d.Id(), "/", 2)
	if len(databaseID) != 2 {
		return diag.Errorf("invalid vkcs_db_database ID: %s", d.Id())
	}

	dbmsID := databaseID[0]
	databaseName := databaseID[1]
	dbmsType := d.Get("dbms_type").(string)

	exists, err := databaseDatabaseExists(DatabaseV1Client, dbmsID, databaseName, dbmsType)
	if err != nil {
		return diag.Errorf("error checking if vkcs_db_database %s exists: %s", d.Id(), err)
	}

	if !exists {
		return nil
	}

	err = databases.Delete(DatabaseV1Client, dbmsID, databaseName, dbmsType).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting vkcs_db_database %s: %s", d.Id(), err)
	}

	return nil
}
