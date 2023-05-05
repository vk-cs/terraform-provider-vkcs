package vkcs

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/datastores"
)

func dataSourceDatabaseDatastoreCapabilities() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatabaseDatastoreCapabilitiesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The `region` to fetch availability zones from, defaults to the provider's `region`.",
			},

			"datastore_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the data store.",
			},

			"datastore_version_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the version of the data store.",
			},

			"capabilities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of data store capability.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of data store capability.",
						},
						"params": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataSourceDatabaseDatastoreCapabilitiesParam(),
						},
						"should_be_on_master": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "This attribute indicates whether a capability applies only to the master node.",
						},
						"allow_major_upgrade": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "This attribute indicates whether a capability can be applied in the next major version of data store.",
						},
						"allow_upgrade_from_backup": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "This attribute indicates whether a capability can be applied to upgrade from backup.",
						},
					},
				},
				Description: "Versions of the datastore.",
			},
		},
		Description: "Use this data source to get capabilities supported for a VKCS datastore. **New since v.0.2.0**.",
	}
}

func dataSourceDatabaseDatastoreCapabilitiesParam() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of a parameter.",
			},
			"required": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Required indicates whether a parameter value must be set.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of value for a parameter.",
			},
			"element_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of element value for a parameter of `list` type.",
			},
			"enum_values": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Supported values for a parameter.",
			},
			"default_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default value for a parameter.",
			},
			"min": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Minimum value for a parameter.",
			},
			"max": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Maximum value for a parameter.",
			},
			"regex": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Regular expression that a parameter value must match.",
			},
			"masked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Masked indicates whether a parameter value must be a boolean mask.",
			},
		},
	}
}

func dataSourceDatabaseDatastoreCapabilitiesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	region := getRegion(d, config)
	dbClient, err := config.DatabaseV1Client(region)
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	dsName := d.Get("datastore_name").(string)
	dsVersionID := d.Get("datastore_version_id").(string)

	capabilities, err := datastores.ListCapabilities(dbClient, dsName, dsVersionID).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_db_backup"))
	}

	flattenedCapabilities := flattenDatabaseDatastoreCapabilities(capabilities)

	d.SetId(fmt.Sprintf("%s/%s/capabilities", dsName, dsVersionID))
	d.Set("capabilities", flattenedCapabilities)

	return nil
}

func flattenDatabaseDatastoreCapabilities(capabilities []datastores.Capability) (r []map[string]interface{}) {
	for _, c := range capabilities {
		r = append(r, map[string]interface{}{
			"name":                      c.Name,
			"description":               c.Description,
			"params":                    flattenDatabaseDatastoreCapabilityParams(c.Params),
			"should_be_on_master":       c.ShouldBeOnMaster,
			"allow_major_upgrade":       c.AllowMajorUpgrade,
			"allow_upgrade_from_backup": c.AllowUpgradeFromBackup,
		})
	}
	return
}

func flattenDatabaseDatastoreCapabilityParams(params map[string]*datastores.CapabilityParam) (r []map[string]interface{}) {
	for name, p := range params {
		var defaultValue string
		switch v := p.DefaultValue.(type) {
		case string:
			defaultValue = v
		case float64:
			defaultValue = strconv.FormatFloat(p.DefaultValue.(float64), 'f', -1, 64)
		}
		r = append(r, map[string]interface{}{
			"name":          name,
			"required":      p.Required,
			"type":          p.Type,
			"element_type":  p.ElementType,
			"enum_values":   p.EnumValues,
			"default_value": defaultValue,
			"min":           p.MinValue,
			"max":           p.MaxValue,
			"regex":         p.Regex,
			"masked":        p.Masked,
		})
	}
	return
}
