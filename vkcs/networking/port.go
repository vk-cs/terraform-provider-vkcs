package networking

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/dns"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/policies"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/utils/terraform/hashcode"
)

type portExtended struct {
	ports.Port
	extradhcpopts.ExtraDHCPOptsExt
	portsecurity.PortSecurityExt
	dns.PortDNSExt
	policies.QoSPolicyExt
	networking.SDNExt
}

func resourceNetworkingPortStateRefreshFunc(client *gophercloud.ServiceClient, portID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := ports.Get(client, portID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return n, "DELETED", nil
			}

			return n, "", err
		}

		return n, n.Status, nil
	}
}

func expandNetworkingPortDHCPOptsCreate(dhcpOpts *schema.Set) []extradhcpopts.CreateExtraDHCPOpt {
	var extraDHCPOpts []extradhcpopts.CreateExtraDHCPOpt

	if dhcpOpts != nil {
		for _, raw := range dhcpOpts.List() {
			rawMap := raw.(map[string]interface{})

			optName := rawMap["name"].(string)
			optValue := rawMap["value"].(string)

			extraDHCPOpts = append(extraDHCPOpts, extradhcpopts.CreateExtraDHCPOpt{
				OptName:   optName,
				OptValue:  optValue,
				IPVersion: gophercloud.IPVersion(4),
			})
		}
	}

	return extraDHCPOpts
}

func expandNetworkingPortDHCPOptsUpdate(oldDHCPopts, newDHCPopts *schema.Set) []extradhcpopts.UpdateExtraDHCPOpt {
	var extraDHCPOpts []extradhcpopts.UpdateExtraDHCPOpt
	var newOptNames []string

	if newDHCPopts != nil {
		for _, raw := range newDHCPopts.List() {
			rawMap := raw.(map[string]interface{})

			optName := rawMap["name"].(string)
			optValue := rawMap["value"].(string)
			// DHCP option name is the primary key, we will check this key below
			newOptNames = append(newOptNames, optName)

			extraDHCPOpts = append(extraDHCPOpts, extradhcpopts.UpdateExtraDHCPOpt{
				OptName:   optName,
				OptValue:  &optValue,
				IPVersion: gophercloud.IPVersion(4),
			})
		}
	}

	if oldDHCPopts != nil {
		for _, raw := range oldDHCPopts.List() {
			rawMap := raw.(map[string]interface{})

			optName := rawMap["name"].(string)

			// if we already add a new option with the same name, it means that we update it, no need to delete
			if !util.StrSliceContains(newOptNames, optName) {
				extraDHCPOpts = append(extraDHCPOpts, extradhcpopts.UpdateExtraDHCPOpt{
					OptName:  optName,
					OptValue: nil,
				})
			}
		}
	}

	return extraDHCPOpts
}

func flattenNetworkingPortDHCPOpts(dhcpOpts extradhcpopts.ExtraDHCPOptsExt) []map[string]interface{} {
	dhcpOptsSet := make([]map[string]interface{}, len(dhcpOpts.ExtraDHCPOpts))

	for i, dhcpOpt := range dhcpOpts.ExtraDHCPOpts {
		dhcpOptsSet[i] = map[string]interface{}{
			"name":  dhcpOpt.OptName,
			"value": dhcpOpt.OptValue,
		}
	}

	return dhcpOptsSet
}

func expandNetworkingPortAllowedAddressPairs(allowedAddressPairs *schema.Set) []ports.AddressPair {
	rawPairs := allowedAddressPairs.List()

	pairs := make([]ports.AddressPair, len(rawPairs))
	for i, raw := range rawPairs {
		rawMap := raw.(map[string]interface{})
		pairs[i] = ports.AddressPair{
			IPAddress:  rawMap["ip_address"].(string),
			MACAddress: rawMap["mac_address"].(string),
		}
	}

	return pairs
}

func flattenNetworkingPortAllowedAddressPairs(mac string, allowedAddressPairs []ports.AddressPair) []map[string]interface{} {
	pairs := make([]map[string]interface{}, len(allowedAddressPairs))

	for i, pair := range allowedAddressPairs {
		pairs[i] = map[string]interface{}{
			"ip_address": pair.IPAddress,
		}
		// Only set the MAC address if it is different than the
		// port's MAC. This means that a specific MAC was set.
		if pair.MACAddress != mac {
			pairs[i]["mac_address"] = pair.MACAddress
		}
	}

	return pairs
}

func expandNetworkingPortFixedIP(d *schema.ResourceData) interface{} {
	// If no_fixed_ip was specified, then just return an empty array.
	// Since no_fixed_ip is mutually exclusive to fixed_ip,
	// we can safely do this.
	//
	// Since we're only concerned about no_fixed_ip being set to "true",
	// GetOk is used.
	if _, ok := d.GetOk("no_fixed_ip"); ok {
		return []interface{}{}
	}

	rawIP := d.Get("fixed_ip").([]interface{})

	if len(rawIP) == 0 {
		return nil
	}

	ip := make([]ports.IP, len(rawIP))
	for i, raw := range rawIP {
		rawMap := raw.(map[string]interface{})
		ip[i] = ports.IP{
			SubnetID:  rawMap["subnet_id"].(string),
			IPAddress: rawMap["ip_address"].(string),
		}
	}
	return ip
}

func resourceNetworkingPortAllowedAddressPairsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-%s", m["ip_address"].(string), m["mac_address"].(string)))

	return hashcode.String(buf.String())
}

func expandNetworkingPortFixedIPToStringSlice(fixedIPs []ports.IP) []string {
	s := make([]string, len(fixedIPs))
	for i, fixedIP := range fixedIPs {
		s[i] = fixedIP.IPAddress
	}

	return s
}
