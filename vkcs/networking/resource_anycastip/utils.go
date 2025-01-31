package resource_anycastip

import "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking/v2/anycastips"

func HasOctavia(associations []anycastips.AnycastIPAssociation) bool {
	for _, association := range associations {
		if association.Type == anycastips.AnycastIPAssociationTypeOctavia {
			return true
		}
	}
	return false
}
