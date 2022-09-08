package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/policies"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

type testPortWithExtensions struct {
	ports.Port
	portsecurity.PortSecurityExt
	policies.QoSPolicyExt
}

func TestAccNetworkingPort_basic(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
				),
			},
		},
	})
}

func TestAccNetworkingPort_noIP(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortNoIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					testAccCheckNetworkingPortCountFixedIPs(&port, 1),
				),
			},
			{
				Config: testAccNetworkingPortNoIPEmptyUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					testAccCheckNetworkingPortCountFixedIPs(&port, 1),
				),
			},
		},
	})
}

func TestAccNetworkingPort_multipleNoIP(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortMultipleNoIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					testAccCheckNetworkingPortCountFixedIPs(&port, 3),
				),
			},
		},
	})
}

func TestAccNetworkingPort_allowedAddressPairs(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var vrrpPort1, vrrpPort2, instancePort ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortAllowedAddressPairs1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.vrrp_subnet", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.vrrp_network", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.vrrp_port_1", &vrrpPort1),
					testAccCheckNetworkingPortExists("vkcs_networking_port.vrrp_port_2", &vrrpPort2),
					testAccCheckNetworkingPortExists("vkcs_networking_port.instance_port", &instancePort),
					testAccCheckNetworkingPortCountAllowedAddressPairs(&instancePort, 2),
					resource.TestCheckResourceAttr("vkcs_networking_port.vrrp_port_1", "description", "test vrrp port"),
				),
			},
			{
				Config: testAccNetworkingPortAllowedAddressPairs2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.vrrp_subnet", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.vrrp_network", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.vrrp_port_1", &vrrpPort1),
					testAccCheckNetworkingPortExists("vkcs_networking_port.vrrp_port_2", &vrrpPort2),
					testAccCheckNetworkingPortExists("vkcs_networking_port.instance_port", &instancePort),
					testAccCheckNetworkingPortCountAllowedAddressPairs(&instancePort, 2),
					resource.TestCheckResourceAttr("vkcs_networking_port.vrrp_port_1", "description", ""),
				),
			},
			{
				Config: testAccNetworkingPortAllowedAddressPairs3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.vrrp_subnet", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.vrrp_network", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.vrrp_port_1", &vrrpPort1),
					testAccCheckNetworkingPortExists("vkcs_networking_port.vrrp_port_2", &vrrpPort2),
					testAccCheckNetworkingPortExists("vkcs_networking_port.instance_port", &instancePort),
					testAccCheckNetworkingPortCountAllowedAddressPairs(&instancePort, 2),
				),
			},
			{
				Config: testAccNetworkingPortAllowedAddressPairs4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.vrrp_subnet", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.vrrp_network", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.vrrp_port_1", &vrrpPort1),
					testAccCheckNetworkingPortExists("vkcs_networking_port.vrrp_port_2", &vrrpPort2),
					testAccCheckNetworkingPortExists("vkcs_networking_port.instance_port", &instancePort),
					testAccCheckNetworkingPortCountAllowedAddressPairs(&instancePort, 1),
				),
			},
			{
				Config: testAccNetworkingPortAllowedAddressPairs5,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.vrrp_subnet", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.vrrp_network", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.vrrp_port_1", &vrrpPort1),
					testAccCheckNetworkingPortExists("vkcs_networking_port.vrrp_port_2", &vrrpPort2),
					testAccCheckNetworkingPortExists("vkcs_networking_port.instance_port", &instancePort),
					testAccCheckNetworkingPortCountAllowedAddressPairs(&instancePort, 0),
				),
			},
		},
	})
}

func TestAccNetworkingPort_allowedAddressPairsNoMAC(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var vrrpPort1, vrrpPort2, instancePort ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortAllowedAddressPairsNoMAC,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.vrrp_subnet", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.vrrp_network", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.vrrp_port_1", &vrrpPort1),
					testAccCheckNetworkingPortExists("vkcs_networking_port.vrrp_port_2", &vrrpPort2),
					testAccCheckNetworkingPortExists("vkcs_networking_port.instance_port", &instancePort),
					testAccCheckNetworkingPortCountAllowedAddressPairs(&instancePort, 2),
				),
			},
		},
	})
}

func TestAccNetworkingPort_multipleFixedIPs(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortMultipleFixedIPs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					testAccCheckNetworkingPortCountFixedIPs(&port, 3),
				),
			},
		},
	})
}

func TestAccNetworkingPort_timeout(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
				),
			},
		},
	})
}

func TestAccNetworkingPort_fixedIPs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortFixedIPs,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "all_fixed_ips.0", "192.168.199.23"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "all_fixed_ips.1", "192.168.199.24"),
				),
			},
		},
	})
}

func TestAccNetworkingPort_updateSecurityGroups(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var secgroup1, secgroup2 groups.SecGroup
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortUpdateSecurityGroups1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_1", &secgroup1),
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_2", &secgroup2),
					testAccCheckNetworkingPortCountSecurityGroups(&port, 1),
				),
			},
			{
				Config: testAccNetworkingPortUpdateSecurityGroups2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_1", &secgroup1),
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_2", &secgroup2),
					testAccCheckNetworkingPortCountSecurityGroups(&port, 1),
				),
			},
			{
				Config: testAccNetworkingPortUpdateSecurityGroups3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_1", &secgroup1),
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_2", &secgroup2),
					testAccCheckNetworkingPortCountSecurityGroups(&port, 2),
				),
			},
			{
				Config: testAccNetworkingPortUpdateSecurityGroups4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_1", &secgroup1),
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_2", &secgroup2),
					testAccCheckNetworkingPortCountSecurityGroups(&port, 1),
				),
			},
			{
				Config: testAccNetworkingPortUpdateSecurityGroups5,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_1", &secgroup1),
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_2", &secgroup2),
					testAccCheckNetworkingPortCountSecurityGroups(&port, 0),
				),
			},
		},
	})
}

func TestAccNetworkingPort_noSecurityGroups(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var secgroup1, secgroup2 groups.SecGroup
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortNoSecurityGroups1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_1", &secgroup1),
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_2", &secgroup2),
					testAccCheckNetworkingPortCountSecurityGroups(&port, 0),
				),
			},
			{
				Config: testAccNetworkingPortNoSecurityGroups2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_1", &secgroup1),
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_2", &secgroup2),
					testAccCheckNetworkingPortCountSecurityGroups(&port, 1),
				),
			},
			{
				Config: testAccNetworkingPortNoSecurityGroups3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_1", &secgroup1),
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_2", &secgroup2),
					testAccCheckNetworkingPortCountSecurityGroups(&port, 2),
				),
			},
			{
				Config: testAccNetworkingPortNoSecurityGroups4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_1", &secgroup1),
					testAccCheckNetworkingSecGroupExists(
						"vkcs_networking_secgroup.secgroup_2", &secgroup2),
					testAccCheckNetworkingPortCountSecurityGroups(&port, 0),
				),
			},
		},
	})
}

func TestAccNetworkingPort_noFixedIP(t *testing.T) {
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortNoFixedIP1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "all_fixed_ips.#", "0"),
				),
			},
			{
				Config: testAccNetworkingPortNoFixedIP2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "all_fixed_ips.#", "1"),
				),
			},
			{
				Config: testAccNetworkingPortNoFixedIP1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "all_fixed_ips.#", "0"),
				),
			},
			{
				Config: testAccNetworkingPortNoFixedIP3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "all_fixed_ips.#", "2"),
				),
			},
			{
				Config: testAccNetworkingPortNoFixedIP1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "all_fixed_ips.#", "0"),
				),
			},
		},
	})
}

func TestAccNetworkingPort_createExtraDHCPOpts(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortCreateExtraDhcpOpts,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "extra_dhcp_option.#", "2"),
				),
			},
		},
	})
}

func TestAccNetworkingPort_updateExtraDHCPOpts(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
				),
			},
			{
				Config: testAccNetworkingPortUpdateExtraDhcpOpts1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "extra_dhcp_option.#", "1"),
				),
			},
			{
				Config: testAccNetworkingPortUpdateExtraDhcpOpts2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "extra_dhcp_option.#", "2"),
				),
			},
			{
				Config: testAccNetworkingPortUpdateExtraDhcpOpts3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "extra_dhcp_option.#", "2"),
				),
			},
			{
				Config: testAccNetworkingPortUpdateExtraDhcpOpts4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "extra_dhcp_option.#", "2"),
				),
			},
			{
				Config: testAccNetworkingPortUpdateExtraDhcpOpts5,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "extra_dhcp_option.#", "2"),
				),
			},
			{
				Config: testAccNetworkingPortUpdateExtraDhcpOpts6,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetExists("vkcs_networking_subnet.subnet_1", &subnet),
					testAccCheckNetworkingNetworkExists("vkcs_networking_network.network_1", &network),
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "extra_dhcp_option.#", "0"),
				),
			},
		},
	})
}

func TestAccNetworkingPort_adminStateUp_omit(t *testing.T) {
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortAdminStateUpOmit,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "admin_state_up", "true"),
					testAccCheckNetworkingPortAdminStateUp(&port, true),
				),
			},
		},
	})
}

func TestAccNetworkingPort_adminStateUp_true(t *testing.T) {
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortAdminStateUpTrue,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "admin_state_up", "true"),
					testAccCheckNetworkingPortAdminStateUp(&port, true),
				),
			},
		},
	})
}

func TestAccNetworkingPort_adminStateUp_false(t *testing.T) {
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortAdminStateUpFalse,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "admin_state_up", "false"),
					testAccCheckNetworkingPortAdminStateUp(&port, false),
				),
			},
		},
	})
}

func TestAccNetworkingPort_adminStateUp_update(t *testing.T) {
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortAdminStateUpOmit,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "admin_state_up", "true"),
					testAccCheckNetworkingPortAdminStateUp(&port, true),
				),
			},
			{
				Config: testAccNetworkingPortAdminStateUpFalse,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "admin_state_up", "false"),
					testAccCheckNetworkingPortAdminStateUp(&port, false),
				),
			},
		},
	})
}

func TestAccNetworkingPort_portSecurity_omit(t *testing.T) {
	var port testPortWithExtensions

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortSecurityOmit,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortWithExtensionsExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "port_security_enabled", "true"),
					testAccCheckNetworkingPortPortSecurityEnabled(&port, true),
				),
			},
			{
				Config: testAccNetworkingPortSecurityDisabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortWithExtensionsExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "port_security_enabled", "false"),
					testAccCheckNetworkingPortPortSecurityEnabled(&port, false),
				),
			},
			{
				Config: testAccNetworkingPortSecurityEnabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortWithExtensionsExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "port_security_enabled", "true"),
					testAccCheckNetworkingPortPortSecurityEnabled(&port, true),
				),
			},
		},
	})
}

func TestAccNetworkingPort_portSecurity_disabled(t *testing.T) {
	var port testPortWithExtensions

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortSecurityDisabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortWithExtensionsExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "port_security_enabled", "false"),
					testAccCheckNetworkingPortPortSecurityEnabled(&port, false),
				),
			},
			{
				Config: testAccNetworkingPortSecurityEnabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortWithExtensionsExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "port_security_enabled", "true"),
					testAccCheckNetworkingPortPortSecurityEnabled(&port, true),
				),
			},
		},
	})
}

func TestAccNetworkingPort_portSecurity_enabled(t *testing.T) {
	var port testPortWithExtensions

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortSecurityEnabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortWithExtensionsExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "port_security_enabled", "true"),
					testAccCheckNetworkingPortPortSecurityEnabled(&port, true),
				),
			},
			{
				Config: testAccNetworkingPortSecurityDisabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingPortWithExtensionsExists("vkcs_networking_port.port_1", &port),
					resource.TestCheckResourceAttr(
						"vkcs_networking_port.port_1", "port_security_enabled", "false"),
					testAccCheckNetworkingPortPortSecurityEnabled(&port, false),
				),
			},
		},
	})
}

func testAccCheckNetworkingPortDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)
	networkingClient, err := config.NetworkingV2Client(osRegionName, defaultSDN)
	if err != nil {
		return fmt.Errorf("Error creating VKCS networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_networking_port" {
			continue
		}

		_, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Port still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingPortExists(n string, port *ports.Port) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(configer)
		networkingClient, err := config.NetworkingV2Client(osRegionName, defaultSDN)
		if err != nil {
			return fmt.Errorf("Error creating VKCS networking client: %s", err)
		}

		found, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Port not found")
		}

		*port = *found

		return nil
	}
}

func testAccCheckNetworkingPortWithExtensionsExists(
	n string, port *testPortWithExtensions) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(configer)
		networkingClient, err := config.NetworkingV2Client(osRegionName, defaultSDN)
		if err != nil {
			return fmt.Errorf("Error creating VKCS networking client: %s", err)
		}

		var p testPortWithExtensions
		err = ports.Get(networkingClient, rs.Primary.ID).ExtractInto(&p)
		if err != nil {
			return err
		}

		if p.ID != rs.Primary.ID {
			return fmt.Errorf("Port not found")
		}

		*port = p

		return nil
	}
}

func testAccCheckNetworkingPortCountFixedIPs(port *ports.Port, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(port.FixedIPs) != expected {
			return fmt.Errorf("Expected %d Fixed IPs, got %d", expected, len(port.FixedIPs))
		}

		return nil
	}
}

func testAccCheckNetworkingPortCountSecurityGroups(port *ports.Port, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(port.SecurityGroups) != expected {
			return fmt.Errorf("Expected %d Security Groups, got %d", expected, len(port.SecurityGroups))
		}

		return nil
	}
}

func testAccCheckNetworkingPortCountAllowedAddressPairs(
	port *ports.Port, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(port.AllowedAddressPairs) != expected {
			return fmt.Errorf("Expected %d Allowed Address Pairs, got %d", expected, len(port.AllowedAddressPairs))
		}

		return nil
	}
}

func testAccCheckNetworkingPortAdminStateUp(port *ports.Port, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if port.AdminStateUp != expected {
			return fmt.Errorf("Port has wrong admin_state_up. Expected %t, got %t", expected, port.AdminStateUp)
		}

		return nil
	}
}

func testAccCheckNetworkingPortPortSecurityEnabled(
	port *testPortWithExtensions, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if port.PortSecurityEnabled != expected {
			return fmt.Errorf("Port has wrong port_security_enabled. Expected %t, got %t", expected, port.PortSecurityEnabled)
		}

		return nil
	}
}

const testAccNetworkingPortBasic = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingPortNoIP = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
  }
}
`

const testAccNetworkingPortNoIPEmptyUpdate = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"
}
`

const testAccNetworkingPortMultipleNoIP = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
  }

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
  }

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
  }
}
`

const testAccNetworkingPortAllowedAddressPairs1 = `
resource "vkcs_networking_network" "vrrp_network" {
  name = "vrrp_network"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "vrrp_subnet" {
  name = "vrrp_subnet"
  cidr = "10.0.0.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  allocation_pool {
    start = "10.0.0.2"
    end = "10.0.0.200"
  }
}

resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_router" "vrrp_router" {
  name = "vrrp_router"
}

resource "vkcs_networking_router_interface" "vrrp_interface" {
  router_id = "${vkcs_networking_router.vrrp_router.id}"
  subnet_id = "${vkcs_networking_subnet.vrrp_subnet.id}"
}

resource "vkcs_networking_port" "vrrp_port_1" {
  name = "vrrp_port_1"
  description = "test vrrp port"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.vrrp_subnet.id}"
    ip_address = "10.0.0.202"
  }
}

resource "vkcs_networking_port" "vrrp_port_2" {
  name = "vrrp_port_2"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.vrrp_subnet.id}"
    ip_address = "10.0.0.201"
  }
}

resource "vkcs_networking_port" "instance_port" {
  name = "instance_port"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  allowed_address_pairs {
    ip_address = "${vkcs_networking_port.vrrp_port_1.fixed_ip.0.ip_address}"
    mac_address = "${vkcs_networking_port.vrrp_port_1.mac_address}"
  }

  allowed_address_pairs {
    ip_address = "${vkcs_networking_port.vrrp_port_2.fixed_ip.0.ip_address}"
    mac_address = "${vkcs_networking_port.vrrp_port_2.mac_address}"
  }
}
`

const testAccNetworkingPortAllowedAddressPairs2 = `
resource "vkcs_networking_network" "vrrp_network" {
  name = "vrrp_network"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "vrrp_subnet" {
  name = "vrrp_subnet"
  cidr = "10.0.0.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  allocation_pool {
    start = "10.0.0.2"
    end = "10.0.0.200"
  }
}

resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_router" "vrrp_router" {
  name = "vrrp_router"
}

resource "vkcs_networking_router_interface" "vrrp_interface" {
  router_id = "${vkcs_networking_router.vrrp_router.id}"
  subnet_id = "${vkcs_networking_subnet.vrrp_subnet.id}"
}

resource "vkcs_networking_port" "vrrp_port_1" {
  name = "vrrp_port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.vrrp_subnet.id}"
    ip_address = "10.0.0.202"
  }
}

resource "vkcs_networking_port" "vrrp_port_2" {
  name = "vrrp_port_2"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.vrrp_subnet.id}"
    ip_address = "10.0.0.201"
  }
}

resource "vkcs_networking_port" "instance_port" {
  name = "instance_port"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  allowed_address_pairs {
    ip_address = "${vkcs_networking_port.vrrp_port_1.fixed_ip.0.ip_address}"
    mac_address = "${vkcs_networking_port.vrrp_port_1.mac_address}"
  }

  allowed_address_pairs {
    ip_address = "${vkcs_networking_port.vrrp_port_2.fixed_ip.0.ip_address}"
    mac_address = "${vkcs_networking_port.vrrp_port_2.mac_address}"
  }
}
`

const testAccNetworkingPortAllowedAddressPairs3 = `
resource "vkcs_networking_network" "vrrp_network" {
  name = "vrrp_network"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "vrrp_subnet" {
  name = "vrrp_subnet"
  cidr = "10.0.0.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  allocation_pool {
    start = "10.0.0.2"
    end = "10.0.0.200"
  }
}

resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_router" "vrrp_router" {
  name = "vrrp_router"
}

resource "vkcs_networking_router_interface" "vrrp_interface" {
  router_id = "${vkcs_networking_router.vrrp_router.id}"
  subnet_id = "${vkcs_networking_subnet.vrrp_subnet.id}"
}

resource "vkcs_networking_port" "vrrp_port_1" {
  name = "vrrp_port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.vrrp_subnet.id}"
    ip_address = "10.0.0.202"
  }
}

resource "vkcs_networking_port" "vrrp_port_2" {
  name = "vrrp_port_2"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.vrrp_subnet.id}"
    ip_address = "10.0.0.201"
  }
}

resource "vkcs_networking_port" "instance_port" {
  name = "instance_port"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.vrrp_network.id}"
  security_group_ids = ["${vkcs_networking_secgroup.secgroup_1.id}"]

  allowed_address_pairs {
    ip_address = "${vkcs_networking_port.vrrp_port_1.fixed_ip.0.ip_address}"
    mac_address = "${vkcs_networking_port.vrrp_port_1.mac_address}"
  }

  allowed_address_pairs {
    ip_address = "${vkcs_networking_port.vrrp_port_2.fixed_ip.0.ip_address}"
    mac_address = "${vkcs_networking_port.vrrp_port_2.mac_address}"
  }
}
`

const testAccNetworkingPortAllowedAddressPairs4 = `
resource "vkcs_networking_network" "vrrp_network" {
  name = "vrrp_network"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "vrrp_subnet" {
  name = "vrrp_subnet"
  cidr = "10.0.0.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  allocation_pool {
    start = "10.0.0.2"
    end = "10.0.0.200"
  }
}

resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_router" "vrrp_router" {
  name = "vrrp_router"
}

resource "vkcs_networking_router_interface" "vrrp_interface" {
  router_id = "${vkcs_networking_router.vrrp_router.id}"
  subnet_id = "${vkcs_networking_subnet.vrrp_subnet.id}"
}

resource "vkcs_networking_port" "vrrp_port_1" {
  name = "vrrp_port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.vrrp_subnet.id}"
    ip_address = "10.0.0.202"
  }
}

resource "vkcs_networking_port" "vrrp_port_2" {
  name = "vrrp_port_2"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.vrrp_subnet.id}"
    ip_address = "10.0.0.201"
  }
}

resource "vkcs_networking_port" "instance_port" {
  name = "instance_port"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.vrrp_network.id}"
  security_group_ids = ["${vkcs_networking_secgroup.secgroup_1.id}"]

  allowed_address_pairs {
    ip_address = "${vkcs_networking_port.vrrp_port_1.fixed_ip.0.ip_address}"
    mac_address = "${vkcs_networking_port.vrrp_port_1.mac_address}"
  }
}
`

const testAccNetworkingPortAllowedAddressPairs5 = `
resource "vkcs_networking_network" "vrrp_network" {
  name = "vrrp_network"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "vrrp_subnet" {
  name = "vrrp_subnet"
  cidr = "10.0.0.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  allocation_pool {
    start = "10.0.0.2"
    end = "10.0.0.200"
  }
}

resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_router" "vrrp_router" {
  name = "vrrp_router"
}

resource "vkcs_networking_router_interface" "vrrp_interface" {
  router_id = "${vkcs_networking_router.vrrp_router.id}"
  subnet_id = "${vkcs_networking_subnet.vrrp_subnet.id}"
}

resource "vkcs_networking_port" "vrrp_port_1" {
  name = "vrrp_port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.vrrp_subnet.id}"
    ip_address = "10.0.0.202"
  }
}

resource "vkcs_networking_port" "vrrp_port_2" {
  name = "vrrp_port_2"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.vrrp_subnet.id}"
    ip_address = "10.0.0.201"
  }
}

resource "vkcs_networking_port" "instance_port" {
  name = "instance_port"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.vrrp_network.id}"
}
`

const testAccNetworkingPortMultipleFixedIPs = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.20"
  }

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.40"
  }
}
`

const testAccNetworkingPortTimeout = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

const testAccNetworkingPortFixedIPs = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.24"
  }

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingPortUpdateSecurityGroups1 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_secgroup" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingPortUpdateSecurityGroups2 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_secgroup" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"
  security_group_ids = ["${vkcs_networking_secgroup.secgroup_1.id}"]

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingPortUpdateSecurityGroups3 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "security_group_1"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_secgroup" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"
  security_group_ids = [
    "${vkcs_networking_secgroup.secgroup_1.id}",
    "${vkcs_networking_secgroup.secgroup_2.id}"
  ]

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingPortUpdateSecurityGroups4 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "security_group"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_secgroup" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"
  security_group_ids = ["${vkcs_networking_secgroup.secgroup_2.id}"]

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingPortUpdateSecurityGroups5 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "security_group"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_secgroup" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"
  security_group_ids = []

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingPortNoSecurityGroups1 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_secgroup" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"
  no_security_groups = true

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingPortNoSecurityGroups2 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_secgroup" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"
  no_security_groups = false
  security_group_ids = ["${vkcs_networking_secgroup.secgroup_1.id}"]

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingPortNoSecurityGroups3 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_secgroup" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"
  no_security_groups = false
  security_group_ids = [
    "${vkcs_networking_secgroup.secgroup_1.id}",
    "${vkcs_networking_secgroup.secgroup_2.id}"
  ]

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingPortNoSecurityGroups4 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_secgroup" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"
  no_security_groups = true

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingPortAllowedAddressPairsNoMAC = `
resource "vkcs_networking_network" "vrrp_network" {
  name = "vrrp_network"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "vrrp_subnet" {
  name = "vrrp_subnet"
  cidr = "10.0.0.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  allocation_pool {
    start = "10.0.0.2"
    end = "10.0.0.200"
  }
}

resource "vkcs_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "vkcs_networking_router" "vrrp_router" {
  name = "vrrp_router"
}

resource "vkcs_networking_router_interface" "vrrp_interface" {
  router_id = "${vkcs_networking_router.vrrp_router.id}"
  subnet_id = "${vkcs_networking_subnet.vrrp_subnet.id}"
}

resource "vkcs_networking_port" "vrrp_port_1" {
  name = "vrrp_port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.vrrp_subnet.id}"
    ip_address = "10.0.0.202"
  }
}

resource "vkcs_networking_port" "vrrp_port_2" {
  name = "vrrp_port_2"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.vrrp_subnet.id}"
    ip_address = "10.0.0.201"
  }
}

resource "vkcs_networking_port" "instance_port" {
  name = "instance_port"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.vrrp_network.id}"

  allowed_address_pairs {
    ip_address = "${vkcs_networking_port.vrrp_port_1.fixed_ip.0.ip_address}"
  }

  allowed_address_pairs {
    ip_address = "${vkcs_networking_port.vrrp_port_2.fixed_ip.0.ip_address}"
  }
}
`

const testAccNetworkingPortNoFixedIP1 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"
  no_fixed_ip = true
}
`

const testAccNetworkingPortNoFixedIP2 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingPortNoFixedIP3 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.24"
  }
}
`

const testAccNetworkingPortCreateExtraDhcpOpts = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  extra_dhcp_option {
    name = "optionA"
    value = "valueA"
  }

  extra_dhcp_option {
    name = "optionB"
    value = "valueB"
  }
}
`

const testAccNetworkingPortUpdateExtraDhcpOpts1 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  extra_dhcp_option {
    name = "optionC"
    value = "valueC"
  }
}
`

const testAccNetworkingPortUpdateExtraDhcpOpts2 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  extra_dhcp_option {
    name = "optionC"
    value = "valueC"
  }

  extra_dhcp_option {
    name = "optionD"
    value = "valueD"
  }
}
`

const testAccNetworkingPortUpdateExtraDhcpOpts3 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  extra_dhcp_option {
    name = "optionD"
    value = "valueD"
  }

  extra_dhcp_option {
    name = "optionE"
    value = "valueE"
  }
}
`

const testAccNetworkingPortUpdateExtraDhcpOpts4 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  extra_dhcp_option {
    name = "optionD"
    value = "valueD"
  }

  extra_dhcp_option {
    name = "optionE"
    value = "valueEE"
  }
}
`

const testAccNetworkingPortUpdateExtraDhcpOpts5 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  extra_dhcp_option {
    name = "optionD"
    value = "valueDD"
  }

  extra_dhcp_option {
    name = "optionE"
    value = "valueEE"
  }
}
`

const testAccNetworkingPortUpdateExtraDhcpOpts6 = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingPortAdminStateUpOmit = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingPortAdminStateUpTrue = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingPortAdminStateUpFalse = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  admin_state_up = "false"
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingPortSecurityOmit = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  no_security_groups = true
  network_id = "${vkcs_networking_network.network_1.id}"

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingPortSecurityDisabled = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  network_id = "${vkcs_networking_network.network_1.id}"
  no_security_groups = true
  port_security_enabled = false

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingPortSecurityEnabled = `
resource "vkcs_networking_network" "network_1" {
  name = "network_1"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_networking_port" "port_1" {
  name = "port_1"
  network_id = "${vkcs_networking_network.network_1.id}"
  no_security_groups = true
  port_security_enabled = true

  fixed_ip {
    subnet_id =  "${vkcs_networking_subnet.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`
