package firewall

import (
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	irules "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/firewall/v2/rules"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

type secgroupRuleExtended struct {
	rules.SecGroupRule
	networking.SDNExt
}

func resourceNetworkingSecGroupRuleStateRefreshFunc(client *gophercloud.ServiceClient, sgRuleID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		sgRule, err := irules.Get(client, sgRuleID).Extract()
		if err != nil {
			if errutil.IsNotFound(err) {
				return sgRule, "DELETED", nil
			}

			return sgRule, "", err
		}

		return sgRule, "ACTIVE", nil
	}
}

func resourceNetworkingSecGroupRuleDirection(direction string) (rules.RuleDirection, error) {
	switch direction {
	case string(rules.DirIngress):
		return rules.DirIngress, nil
	case string(rules.DirEgress):
		return rules.DirEgress, nil
	}

	return "", fmt.Errorf("unknown direction for vkcs_networking_secgroup_rule: %s", direction)
}

func resourceNetworkingSecGroupRuleEtherType(etherType string) (rules.RuleEtherType, error) {
	switch etherType {
	case string(rules.EtherType4):
		return rules.EtherType4, nil
	case string(rules.EtherType6):
		return rules.EtherType6, nil
	}

	return "", fmt.Errorf("unknown ether type for vkcs_networking_secgroup_rule: %s", etherType)
}

func resourceNetworkingSecGroupRuleProtocol(protocol string) (rules.RuleProtocol, error) {
	switch protocol {
	case string(rules.ProtocolAH):
		return rules.ProtocolAH, nil
	case string(rules.ProtocolDCCP):
		return rules.ProtocolDCCP, nil
	case string(rules.ProtocolEGP):
		return rules.ProtocolEGP, nil
	case string(rules.ProtocolESP):
		return rules.ProtocolESP, nil
	case string(rules.ProtocolGRE):
		return rules.ProtocolGRE, nil
	case string(rules.ProtocolICMP):
		return rules.ProtocolICMP, nil
	case string(rules.ProtocolIGMP):
		return rules.ProtocolIGMP, nil
	case string(rules.ProtocolOSPF):
		return rules.ProtocolOSPF, nil
	case string(rules.ProtocolPGM):
		return rules.ProtocolPGM, nil
	case string(rules.ProtocolRSVP):
		return rules.ProtocolRSVP, nil
	case string(rules.ProtocolSCTP):
		return rules.ProtocolSCTP, nil
	case string(rules.ProtocolTCP):
		return rules.ProtocolTCP, nil
	case string(rules.ProtocolUDP):
		return rules.ProtocolUDP, nil
	case string(rules.ProtocolUDPLite):
		return rules.ProtocolUDPLite, nil
	case string(rules.ProtocolVRRP):
		return rules.ProtocolVRRP, nil
	}

	// If the protocol wasn't matched above, see if it's an integer.
	_, err := strconv.Atoi(protocol)
	if err == nil {
		return rules.RuleProtocol(protocol), nil
	}

	return "", fmt.Errorf("unknown protocol for vkcs_networking_secgroup_rule: %s", protocol)
}

var secgroupRuleProtocolNumbers = map[string]string{
	string(rules.ProtocolICMP):    "1",
	string(rules.ProtocolIGMP):    "2",
	string(rules.ProtocolTCP):     "6",
	string(rules.ProtocolEGP):     "8",
	string(rules.ProtocolUDP):     "17",
	string(rules.ProtocolDCCP):    "33",
	string(rules.ProtocolRSVP):    "46",
	string(rules.ProtocolGRE):     "47",
	string(rules.ProtocolESP):     "50",
	string(rules.ProtocolAH):      "51",
	string(rules.ProtocolOSPF):    "89",
	string(rules.ProtocolVRRP):    "112",
	string(rules.ProtocolPGM):     "113",
	string(rules.ProtocolSCTP):    "132",
	string(rules.ProtocolUDPLite): "136",
}

func secgroupRuleProtocolNumber(protocol string) string {
	if number, ok := secgroupRuleProtocolNumbers[protocol]; ok {
		return number
	}

	return protocol
}

// Sprut returns a protocol name even if the rule was created with a protocol number ("6" -> "tcp").
func suppressSecGroupRuleProtocolDiffs(_, old, new string, _ *schema.ResourceData) bool {
	return secgroupRuleProtocolNumber(old) == secgroupRuleProtocolNumber(new)
}
