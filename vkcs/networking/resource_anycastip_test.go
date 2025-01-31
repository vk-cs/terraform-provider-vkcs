package networking_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccNetworkingAnycastIPResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccNetworkingAnycastIPResourceBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_networking_anycastip.basic", "name", "tfacc-anycastip-basic"),
					resource.TestCheckResourceAttr("vkcs_networking_anycastip.basic", "description", "tfacc-anycastip-basic-description"),
				),
			},
			acctest.ImportStep("vkcs_networking_anycastip.basic"),
		},
	})
}

func TestAccNetworkingAnycastIPResource_loadbalancer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccNetworkingAnycastIPResourceLoadbalancer),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_networking_anycastip.loadbalancer", "name", "tfacc-anycastip-loadbalancer"),
					resource.TestCheckResourceAttr("vkcs_networking_anycastip.loadbalancer", "description", "tfacc-anycastip-loadbalancer-description"),
					resource.TestCheckResourceAttr("vkcs_networking_anycastip.loadbalancer", "associations.#", "2"),
					resource.TestCheckResourceAttr("vkcs_networking_anycastip.loadbalancer", "health_check.type", "ICMP"),
				),
			},
			acctest.ImportStep("vkcs_networking_anycastip.loadbalancer"),
		},
	})
}

func TestAccNetworkingAnycastIPResource_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccNetworkingAnycastIPResourceUpdateOld),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_networking_anycastip.basic", "name", "tfacc-anycastip-update-old"),
					resource.TestCheckResourceAttr("vkcs_networking_anycastip.basic", "description", "tfacc-anycastip-update-old-description"),
					resource.TestCheckResourceAttr("vkcs_networking_anycastip.basic", "health_check.type", "TCP"),
					resource.TestCheckResourceAttr("vkcs_networking_anycastip.basic", "health_check.port", "1337"),
				),
			},
			acctest.ImportStep("vkcs_networking_anycastip.basic"),
			{
				Config: acctest.AccTestRenderConfig(testAccNetworkingAnycastIPResourceUpdateNew),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_networking_anycastip.basic", "name", "tfacc-anycastip-update-new"),
					resource.TestCheckResourceAttr("vkcs_networking_anycastip.basic", "description", "tfacc-anycastip-update-new-description"),
					resource.TestCheckResourceAttr("vkcs_networking_anycastip.basic", "health_check.type", "TCP"),
					resource.TestCheckResourceAttr("vkcs_networking_anycastip.basic", "health_check.port", "1338"),
				),
			},
			acctest.ImportStep("vkcs_networking_anycastip.basic"),
		},
	})
}

const testAccNetworkingAnycastIPResourceBasic = `
{{.BaseExtNetwork}}

resource "vkcs_networking_anycastip" "basic" {
  name        = "tfacc-anycastip-basic"
  description = "tfacc-anycastip-basic-description"

  network_id = data.vkcs_networking_network.extnet.id
}
`

const testAccNetworkingAnycastIPResourceLoadbalancer = `
{{.BaseExtNetwork}}

resource "vkcs_networking_network" "app" {
  name        = "app-tf-example-sprut"
  description = "Application network"
  sdn         = "sprut"
}

resource "vkcs_networking_subnet" "app" {
  name       = "app-tf-example-sprut"
  network_id = vkcs_networking_network.app.id
  cidr       = "192.168.199.0/24"
  sdn        = "sprut"
}

resource "vkcs_networking_router" "router" {
  name                = "router-tf-example"
  external_network_id = data.vkcs_networking_network.extnet.id
  sdn                 = "sprut"
}

resource "vkcs_networking_router_interface" "app" {
  router_id = vkcs_networking_router.router.id
  subnet_id = vkcs_networking_subnet.app.id
  sdn       = "sprut"
}

resource "vkcs_lb_loadbalancer" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = vkcs_networking_subnet.app.id

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "vkcs_lb_loadbalancer" "loadbalancer_2" {
  name          = "loadbalancer_2"
  vip_subnet_id = vkcs_networking_subnet.app.id

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "vkcs_networking_anycastip" "loadbalancer" {
  name        = "tfacc-anycastip-loadbalancer"
  description = "tfacc-anycastip-loadbalancer-description"

  network_id = data.vkcs_networking_network.extnet.id
  associations = [
    {
      id   = vkcs_lb_loadbalancer.loadbalancer_1.vip_port_id
      type = "octavia"
    },
    {
      id   = vkcs_lb_loadbalancer.loadbalancer_2.vip_port_id
      type = "octavia"
    }
  ]
}
`

const testAccNetworkingAnycastIPResourceUpdateOld = `
{{.BaseExtNetwork}}

resource "vkcs_networking_anycastip" "basic" {
  name        = "tfacc-anycastip-update-old"
  description = "tfacc-anycastip-update-old-description"

  network_id = data.vkcs_networking_network.extnet.id
  health_check = {
    type = "TCP"
    port = 1337
  }
}
`

const testAccNetworkingAnycastIPResourceUpdateNew = `
{{.BaseExtNetwork}}

resource "vkcs_networking_anycastip" "basic" {
  name        = "tfacc-anycastip-update-new"
  description = "tfacc-anycastip-update-new-description"

  network_id = data.vkcs_networking_network.extnet.id
  health_check = {
    type = "TCP"
    port = 1338
  }
}
`
