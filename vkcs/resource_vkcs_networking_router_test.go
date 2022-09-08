package vkcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
)

func TestAccNetworkingRouter_basic(t *testing.T) {
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingRouterBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingRouterExists("vkcs_networking_router.router_1", &router),
					resource.TestCheckResourceAttr(
						"vkcs_networking_router.router_1", "description", "router description"),
				),
			},
			{
				Config: testAccNetworkingRouterUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_networking_router.router_1", "name", "router_2"),
					resource.TestCheckResourceAttr(
						"vkcs_networking_router.router_1", "description", ""),
				),
			},
		},
	})
}

func TestAccNetworkingRouter_updateExternalGateway(t *testing.T) {
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingRouterUpdateExternalGateway1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingRouterExists("vkcs_networking_router.router_1", &router),
				),
			},
			{
				Config: testAccNetworkingRouterUpdateExternalGateway2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vkcs_networking_router.router_1", "external_network_id", osExtGwID),
				),
			},
		},
	})
}

func TestAccNetworkingRouter_vendor_opts(t *testing.T) {
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingRouterVendorOpts(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingRouterExists("vkcs_networking_router.router_1", &router),
					resource.TestCheckResourceAttr(
						"vkcs_networking_router.router_1", "external_network_id", osExtGwID),
				),
			},
		},
	})
}

func testAccCheckNetworkingRouterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(configer)
	networkingClient, err := config.NetworkingV2Client(osRegionName, defaultSDN)
	if err != nil {
		return fmt.Errorf("Error creating VKCS networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vkcs_networking_router" {
			continue
		}

		_, err := routers.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Router still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingRouterExists(n string, router *routers.Router) resource.TestCheckFunc {
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

		found, err := routers.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Router not found")
		}

		*router = *found

		return nil
	}
}

const testAccNetworkingRouterBasic = `
resource "vkcs_networking_router" "router_1" {
  name = "router_1"
  description = "router description"
  admin_state_up = "true"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

const testAccNetworkingRouterUpdate = `
resource "vkcs_networking_router" "router_1" {
  name = "router_2"
  admin_state_up = "true"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

func testAccNetworkingRouterVendorOpts() string {
	return fmt.Sprintf(`
%s

resource "vkcs_networking_router" "router_1" {
  name = "router_1"
  admin_state_up = "true"
  external_network_id = data.vkcs_networking_network.extnet.id
  vendor_options {
    set_router_gateway_after_create = true
  }
}
`, testAccBaseExtNetwork)
}

const testAccNetworkingRouterUpdateExternalGateway1 = `
resource "vkcs_networking_router" "router_1" {
  name = "router"
  admin_state_up = "true"
}
`

func testAccNetworkingRouterUpdateExternalGateway2() string {
	return fmt.Sprintf(`
%s

resource "vkcs_networking_router" "router_1" {
  name = "router"
  admin_state_up = "true"
  external_network_id = data.vkcs_networking_network.extnet.id
}
`, testAccBaseExtNetwork)
}
