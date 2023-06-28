package vpnaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/siteconnections"
)

func TestAccVPNaaSSiteConnection_basic(t *testing.T) {
	var conn siteconnections.Connection
	resource.Test(t, resource.TestCase{
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckSiteConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccSiteConnectionBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteConnectionExists("vkcs_vpnaas_site_connection.conn_1", &conn),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_site_connection.conn_1", "ikepolicy_id", &conn.IKEPolicyID),
					resource.TestCheckResourceAttr("vkcs_vpnaas_site_connection.conn_1", "admin_state_up", "true"),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_site_connection.conn_1", "psk", &conn.PSK),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_site_connection.conn_1", "ipsecpolicy_id", &conn.IPSecPolicyID),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_site_connection.conn_1", "vpnservice_id", &conn.VPNServiceID),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_site_connection.conn_1", "local_ep_group_id", &conn.LocalEPGroupID),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_site_connection.conn_1", "local_id", &conn.LocalID),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_site_connection.conn_1", "peer_ep_group_id", &conn.PeerEPGroupID),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_site_connection.conn_1", "name", &conn.Name),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_site_connection.conn_1", "dpd.0.action", &conn.DPD.Action),
					resource.TestCheckResourceAttr("vkcs_vpnaas_site_connection.conn_1", "dpd.0.timeout", "42"),
					resource.TestCheckResourceAttr("vkcs_vpnaas_site_connection.conn_1", "dpd.0.interval", "21"),
				),
			},
		},
	})
}

func testAccCheckSiteConnectionDestroy(s *terraform.State) error {
	config := acctest.AccTestProvider.Meta().(clients.Config)
	networkingClient, err := config.NetworkingV2Client(acctest.OsRegionName, networking.DefaultSDN)
	if err != nil {
		return fmt.Errorf("Error creating VKCS networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_vpnaas_site_connection" {
			continue
		}
		_, err = siteconnections.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Site connection (%s) still exists", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckSiteConnectionExists(n string, conn *siteconnections.Connection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.AccTestProvider.Meta().(clients.Config)
		networkingClient, err := config.NetworkingV2Client(acctest.OsRegionName, networking.DefaultSDN)
		if err != nil {
			return fmt.Errorf("Error creating VKCS networking client: %s", err)
		}

		var found *siteconnections.Connection

		found, err = siteconnections.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*conn = *found

		return nil
	}
}

const testAccSiteConnectionBasic = `
	{{.BaseExtNetwork}}
	
	resource "vkcs_networking_network" "network_1" {
		name           = "tf_test_network"
  		admin_state_up = "true"
	}

	resource "vkcs_networking_subnet" "subnet_1" {
  		network_id = vkcs_networking_network.network_1.id
  		cidr       = "192.168.199.0/24"
	}

	resource "vkcs_networking_router" "router_1" {
  		name             = "my_router"
  		external_network_id = data.vkcs_networking_network.extnet.id
	}

	resource "vkcs_networking_router_interface" "router_interface_1" {
  		router_id = vkcs_networking_router.router_1.id
  		subnet_id = vkcs_networking_subnet.subnet_1.id
	}

	resource "vkcs_vpnaas_service" "service_1" {
		router_id = vkcs_networking_router.router_1.id
		admin_state_up = "false"
	}

	resource "vkcs_vpnaas_ipsec_policy" "policy_1" {
	}

	resource "vkcs_vpnaas_ike_policy" "policy_2" {
	}

	resource "vkcs_vpnaas_endpoint_group" "group_1" {
		type = "cidr"
		endpoints = ["10.0.0.24/24", "10.0.0.25/24"]
	}
	resource "vkcs_vpnaas_endpoint_group" "group_2" {
		type = "subnet"
		endpoints = [ vkcs_networking_subnet.subnet_1.id ]
	}

	resource "vkcs_vpnaas_site_connection" "conn_1" {
		name = "connection_1"
		ikepolicy_id = vkcs_vpnaas_ike_policy.policy_2.id
		ipsecpolicy_id = vkcs_vpnaas_ipsec_policy.policy_1.id
		vpnservice_id = vkcs_vpnaas_service.service_1.id
		psk = "secret"
		peer_address = "192.168.10.1"
		peer_id = "192.168.10.1"
		local_ep_group_id = vkcs_vpnaas_endpoint_group.group_2.id
		peer_ep_group_id = vkcs_vpnaas_endpoint_group.group_1.id
		dpd {
			action   = "restart"
			timeout  = 42
			interval = 21
		}
		depends_on = ["vkcs_networking_router_interface.router_interface_1"]
	}
	`
