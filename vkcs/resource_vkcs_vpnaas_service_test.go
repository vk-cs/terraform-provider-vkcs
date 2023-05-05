package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/services"
)

func TestAccVPNaaSService_basic(t *testing.T) {
	var service services.Service
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRenderConfig(testAccServiceBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceExists("vkcs_vpnaas_service.service_1", &service),
					resource.TestCheckResourceAttrPtr("vkcs_vpnaas_service.service_1", "router_id", &service.RouterID),
					resource.TestCheckResourceAttr("vkcs_vpnaas_service.service_1", "admin_state_up", "false"),
				),
			},
		},
	})
}

func testAccCheckServiceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(clients.Config)
	networkingClient, err := config.NetworkingV2Client(osRegionName, networking.DefaultSDN)
	if err != nil {
		return fmt.Errorf("Error creating VKCS networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_vpnaas_service" {
			continue
		}
		_, err = services.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Service (%s) still exists", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckServiceExists(n string, serv *services.Service) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(clients.Config)
		networkingClient, err := config.NetworkingV2Client(osRegionName, networking.DefaultSDN)
		if err != nil {
			return fmt.Errorf("Error creating VKCS networking client: %s", err)
		}

		var found *services.Service

		found, err = services.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*serv = *found

		return nil
	}
}

const testAccServiceBasic = `
{{.BaseExtNetwork}}

	resource "vkcs_networking_router" "router_1" {
	  name = "router_1"
	  admin_state_up = "true"
	  external_network_id = data.vkcs_networking_network.extnet.id
	}

	resource "vkcs_vpnaas_service" "service_1" {
		router_id = vkcs_networking_router.router_1.id
		admin_state_up = "false"
	}
	`
