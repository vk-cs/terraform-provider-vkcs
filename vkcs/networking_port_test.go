package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

func TestExpandNetworkingPortDHCPOptsCreate(t *testing.T) {
	r := resourceNetworkingPort()
	d := r.TestResourceData()
	d.SetId("1")
	dhcpOpts1 := map[string]interface{}{
		"name":  "A",
		"value": "true",
	}
	dhcpOpts2 := map[string]interface{}{
		"name":  "B",
		"value": "false",
	}
	extraDHCPOpts := []map[string]interface{}{dhcpOpts1, dhcpOpts2}
	d.Set("extra_dhcp_option", extraDHCPOpts)

	expectedDHCPOptions := []extradhcpopts.CreateExtraDHCPOpt{
		{
			OptName:   "B",
			OptValue:  "false",
			IPVersion: gophercloud.IPVersion(4),
		},
		{
			OptName:   "A",
			OptValue:  "true",
			IPVersion: gophercloud.IPVersion(4),
		},
	}

	actualDHCPOptions := expandNetworkingPortDHCPOptsCreate(d.Get("extra_dhcp_option").(*schema.Set))

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestExpandNetworkingPortDHCPOptsEmptyCreate(t *testing.T) {
	r := resourceNetworkingPort()
	d := r.TestResourceData()
	d.SetId("1")

	expectedDHCPOptions := []extradhcpopts.CreateExtraDHCPOpt{}

	actualDHCPOptions := expandNetworkingPortDHCPOptsCreate(d.Get("extra_dhcp_option").(*schema.Set))

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestExpandNetworkingPortDHCPOptsUpdate(t *testing.T) {
	r := resourceNetworkingPort()
	d := r.TestResourceData()
	d.SetId("1")
	dhcpOpts1 := map[string]interface{}{
		"name":  "A",
		"value": "true",
	}
	dhcpOpts2 := map[string]interface{}{
		"name":  "B",
		"value": "false",
	}
	extraDHCPOpts := []map[string]interface{}{dhcpOpts1, dhcpOpts2}
	d.Set("extra_dhcp_option", extraDHCPOpts)

	optsValueTrue := "true"
	optsValueFalse := "false"
	expectedDHCPOptions := []extradhcpopts.UpdateExtraDHCPOpt{
		{
			OptName:   "B",
			OptValue:  &optsValueFalse,
			IPVersion: gophercloud.IPVersion(4),
		},
		{
			OptName:   "A",
			OptValue:  &optsValueTrue,
			IPVersion: gophercloud.IPVersion(4),
		},
	}

	actualDHCPOptions := expandNetworkingPortDHCPOptsUpdate(nil, d.Get("extra_dhcp_option").(*schema.Set))

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestExpandNetworkingPortDHCPOptsEmptyUpdate(t *testing.T) {
	r := resourceNetworkingPort()
	d := r.TestResourceData()
	d.SetId("1")

	var expectedDHCPOptions []extradhcpopts.UpdateExtraDHCPOpt

	actualDHCPOptions := expandNetworkingPortDHCPOptsUpdate(nil, d.Get("extra_dhcp_option").(*schema.Set))

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestExpandNetworkingPortDHCPOptsDelete(t *testing.T) {
	r := resourceNetworkingPort()
	d := r.TestResourceData()
	d.SetId("1")
	dhcpOpts1 := map[string]interface{}{
		"name":  "A",
		"value": "true",
	}
	dhcpOpts2 := map[string]interface{}{
		"name":  "B",
		"value": "false",
	}
	extraDHCPOpts := []map[string]interface{}{dhcpOpts1, dhcpOpts2}
	d.Set("extra_dhcp_option", extraDHCPOpts)

	expectedDHCPOptions := []extradhcpopts.UpdateExtraDHCPOpt{
		{
			OptName: "B",
		},
		{
			OptName: "A",
		},
	}

	actualDHCPOptions := expandNetworkingPortDHCPOptsUpdate(d.Get("extra_dhcp_option").(*schema.Set), nil)

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestFlattenNetworkingPort2DHCPOptions(t *testing.T) {
	dhcpOptions := extradhcpopts.ExtraDHCPOptsExt{
		ExtraDHCPOpts: []extradhcpopts.ExtraDHCPOpt{
			{
				OptName:   "A",
				OptValue:  "true",
				IPVersion: 4,
			},
			{
				OptName:   "B",
				OptValue:  "false",
				IPVersion: 6,
			},
		},
	}

	expectedDHCPOptions := []map[string]interface{}{
		{
			"name":  "A",
			"value": "true",
		},
		{
			"name":  "B",
			"value": "false",
		},
	}

	actualDHCPOptions := flattenNetworkingPortDHCPOpts(dhcpOptions)

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestExpandNetworkingPortAllowedAddressPairs(t *testing.T) {
	r := resourceNetworkingPort()
	d := r.TestResourceData()
	d.SetId("1")
	addressPairs1 := map[string]interface{}{
		"ip_address":  "192.0.2.1",
		"mac_address": "mac1",
	}
	addressPairs2 := map[string]interface{}{
		"ip_address":  "198.51.100.1",
		"mac_address": "mac2",
	}
	allowedAddressPairs := []map[string]interface{}{addressPairs1, addressPairs2}
	d.Set("allowed_address_pairs", allowedAddressPairs)

	expectedAllowedAddressPairs := []ports.AddressPair{
		{
			IPAddress:  "192.0.2.1",
			MACAddress: "mac1",
		},
		{
			IPAddress:  "198.51.100.1",
			MACAddress: "mac2",
		},
	}

	actualAllowedAddressPairs := expandNetworkingPortAllowedAddressPairs(d.Get("allowed_address_pairs").(*schema.Set))

	assert.ElementsMatch(t, expectedAllowedAddressPairs, actualAllowedAddressPairs)
}

func TestFlattenNetworkingPortAllowedAddressPairs(t *testing.T) {
	allowedAddressPairs := []ports.AddressPair{
		{
			IPAddress:  "192.0.2.1",
			MACAddress: "mac1",
		},
		{
			IPAddress:  "198.51.100.1",
			MACAddress: "mac2",
		},
	}
	mac := "mac3"

	expectedAllowedAddressPairs := []map[string]interface{}{
		{
			"ip_address":  "192.0.2.1",
			"mac_address": "mac1",
		},
		{
			"ip_address":  "198.51.100.1",
			"mac_address": "mac2",
		},
	}

	actualAllowedAddressPairs := flattenNetworkingPortAllowedAddressPairs(mac, allowedAddressPairs)

	assert.ElementsMatch(t, expectedAllowedAddressPairs, actualAllowedAddressPairs)
}

func TestExpandNetworkingPortFixedIPNoFixedIPs(t *testing.T) {
	r := resourceNetworkingPort()
	d := r.TestResourceData()
	d.SetId("1")
	d.Set("no_fixed_ip", true)

	actualFixedIP := expandNetworkingPortFixedIP(d)

	assert.Empty(t, actualFixedIP)
}

func TestExpandNetworkingPortFixedIPSomeFixedIPs(t *testing.T) {
	r := resourceNetworkingPort()
	d := r.TestResourceData()
	d.SetId("1")
	fixedIP1 := map[string]interface{}{
		"subnet_id":  "aaa",
		"ip_address": "192.0.201.101",
	}
	fixedIP2 := map[string]interface{}{
		"subnet_id":  "bbb",
		"ip_address": "192.0.202.102",
	}
	fixedIP := []map[string]interface{}{fixedIP1, fixedIP2}
	d.Set("fixed_ip", fixedIP)

	expectedFixedIP := []ports.IP{
		{
			SubnetID:  "aaa",
			IPAddress: "192.0.201.101",
		},
		{
			SubnetID:  "bbb",
			IPAddress: "192.0.202.102",
		},
	}

	actualFixedIP := expandNetworkingPortFixedIP(d)

	assert.ElementsMatch(t, expectedFixedIP, actualFixedIP)
}
