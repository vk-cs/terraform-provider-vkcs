package firewall

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
)

func resourceNetworkingSecGroupRuleStateRefreshFunc(client *gophercloud.ServiceClient, sgRuleID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		sgRule, err := rules.Get(client, sgRuleID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
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
