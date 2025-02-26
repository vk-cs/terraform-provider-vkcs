package vpnaas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/services"
	iservices "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/vpnaas/v2/services"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func TestAccVPNaaSService_basic(t *testing.T) {
	var service services.Service
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.AccTestPreCheck(t) },
		ProviderFactories: acctest.AccTestProviders,
		CheckDestroy:      testAccCheckServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccServiceBasic),
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
	config := acctest.AccTestProvider.Meta().(clients.Config)
	networkingClient, err := config.NetworkingV2Client(acctest.OsRegionName, networking.NeutronSDN)
	if err != nil {
		return fmt.Errorf("Error creating VKCS networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_vpnaas_service" {
			continue
		}
		_, err = iservices.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Service (%s) still exists", rs.Primary.ID)
		}
		if !errutil.Is(err, 404) {
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

		config := acctest.AccTestProvider.Meta().(clients.Config)
		networkingClient, err := config.NetworkingV2Client(acctest.OsRegionName, networking.NeutronSDN)
		if err != nil {
			return fmt.Errorf("Error creating VKCS networking client: %s", err)
		}

		var found *services.Service

		found, err = iservices.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*serv = *found

		return nil
	}
}

const testAccServiceBasic = `
{{.BaseExtNetworkNeutron}}

	resource "vkcs_networking_router" "router_1" {
	  name = "router_1"
	  admin_state_up = "true"
	  external_network_id = data.vkcs_networking_network.extnet.id
	  sdn = "neutron"
	}

	resource "vkcs_vpnaas_service" "service_1" {
		router_id = vkcs_networking_router.router_1.id
		admin_state_up = "false"
		sdn = "neutron"
	}
	`
