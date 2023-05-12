package networking

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

// networkingSubnetStateRefreshFunc returns a standard resource.StateRefreshFunc to wait for subnet status.
func networkingSubnetStateRefreshFunc(client *gophercloud.ServiceClient, subnetID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		subnet, err := subnets.Get(client, subnetID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return subnet, "DELETED", nil
			}

			return nil, "", err
		}

		return subnet, "ACTIVE", nil
	}
}

// networkingSubnetStateRefreshFuncDelete returns a special case resource.StateRefreshFunc to try to delete a subnet.
func networkingSubnetStateRefreshFuncDelete(networkingClient *gophercloud.ServiceClient, subnetID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete vkcs_networking_subnet %s", subnetID)

		s, err := subnets.Get(networkingClient, subnetID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted vkcs_networking_subnet %s", subnetID)
				return s, "DELETED", nil
			}

			return s, "ACTIVE", err
		}

		err = subnets.Delete(networkingClient, subnetID).ExtractErr()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted vkcs_networking_subnet %s", subnetID)
				return s, "DELETED", nil
			}
			// Subnet is still in use - we can retry.
			if _, ok := err.(gophercloud.ErrDefault409); ok {
				return s, "ACTIVE", nil
			}

			return s, "ACTIVE", err
		}

		log.Printf("[DEBUG] vkcs_networking_subnet %s is still active", subnetID)

		return s, "ACTIVE", nil
	}
}

// networkingSubnetGetRawAllocationPoolsValueToExpand selects the resource argument to populate
// the allocations pool value.
func networkingSubnetGetRawAllocationPoolsValueToExpand(d *schema.ResourceData) []interface{} {
	result := d.Get("allocation_pool").(*schema.Set).List()
	return result
}

// expandNetworkingSubnetAllocationPools returns a slice of subnets.AllocationPool structs.
func expandNetworkingSubnetAllocationPools(allocationPools []interface{}) []subnets.AllocationPool {
	result := make([]subnets.AllocationPool, len(allocationPools))
	for i, raw := range allocationPools {
		rawMap := raw.(map[string]interface{})

		result[i] = subnets.AllocationPool{
			Start: rawMap["start"].(string),
			End:   rawMap["end"].(string),
		}
	}

	return result
}

// flattenNetworkingSubnetAllocationPools allows to flatten slice of subnets.AllocationPool structs into
// a slice of maps.
func flattenNetworkingSubnetAllocationPools(allocationPools []subnets.AllocationPool) []map[string]interface{} {
	result := make([]map[string]interface{}, len(allocationPools))
	for i, allocationPool := range allocationPools {
		pool := make(map[string]interface{})
		pool["start"] = allocationPool.Start
		pool["end"] = allocationPool.End

		result[i] = pool
	}

	return result
}

func networkingSubnetAllocationPoolsMatch(oldPools, newPools []interface{}) bool {
	if len(oldPools) != len(newPools) {
		return false
	}

	for _, newPool := range newPools {
		var found bool

		newPoolPool := newPool.(map[string]interface{})
		newStart := newPoolPool["start"].(string)
		newEnd := newPoolPool["end"].(string)

		for _, oldPool := range oldPools {
			oldPoolPool := oldPool.(map[string]interface{})
			oldStart := oldPoolPool["start"].(string)
			oldEnd := oldPoolPool["end"].(string)

			if oldStart == newStart && oldEnd == newEnd {
				found = true
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func networkingSubnetDNSNameserverAreUnique(raw []interface{}) error {
	set := make(map[string]struct{})
	for _, rawNS := range raw {
		nameserver, ok := rawNS.(string)
		if ok {
			if _, exists := set[nameserver]; exists {
				return fmt.Errorf("got duplicate nameserver %s", nameserver)
			}
			set[nameserver] = struct{}{}
		}
	}

	return nil
}
