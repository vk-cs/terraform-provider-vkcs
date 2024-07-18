package networking

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func NetworkingReadAttributesTags(d *schema.ResourceData, tags []string) {
	util.ExpandObjectReadTags(d, tags)
}

func NetworkingV2UpdateAttributesTags(d *schema.ResourceData) []string {
	return util.ExpandObjectUpdateTags(d)
}

func NetworkingAttributesTags(d *schema.ResourceData) []string {
	return util.ExpandObjectTags(d)
}

func GetSDN(d *schema.ResourceData) string {
	if v, ok := d.GetOk("sdn"); ok {
		return v.(string)
	}

	return networking.SearchInAllSDNs
}

func ValidateSDN() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{"neutron", "sprut"}, true))
}
