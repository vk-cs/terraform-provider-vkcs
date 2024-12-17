package networking_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccNetworkingPortDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingPortDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.vkcs_networking_port.port_1", "id",
						"vkcs_networking_port.port_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.vkcs_networking_port.port_2", "id",
						"vkcs_networking_port.port_2", "id"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_port.port_2", "allowed_address_pairs.#", "2"),
					resource.TestCheckResourceAttrPair(
						"data.vkcs_networking_port.port_3", "id",
						"vkcs_networking_port.port_1", "id"),
					resource.TestCheckResourceAttr(
						"data.vkcs_networking_port.port_3", "all_fixed_ips.#", "2"),
				),
			},
		},
	})
}

const testAccNetworkingPortDataSourceBasic = `
resource "vkcs_networking_network" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name       = "subnet_1"
  network_id = vkcs_networking_network.network_1.id
  cidr       = "10.0.0.0/24"
}

data "vkcs_networking_secgroup" "default" {
  name = "default"
  sdn = "neutron"
}

resource "vkcs_networking_port" "port_1" {
  name           = "port"
  description    = "test port"
  network_id     = vkcs_networking_network.network_1.id
  admin_state_up = true

  security_group_ids = [
    data.vkcs_networking_secgroup.default.id,
  ]

  fixed_ip {
    subnet_id = vkcs_networking_subnet.subnet_1.id
  }

  fixed_ip {
    subnet_id = vkcs_networking_subnet.subnet_1.id
  }
}

resource "vkcs_networking_port" "port_2" {
  name               = "port"
  description        = "test port"
  network_id         = vkcs_networking_network.network_1.id
  admin_state_up     = true
  no_security_groups = true

  tags = [
    "foo",
    "bar",
  ]

  allowed_address_pairs {
    ip_address  = "10.0.0.201"
    mac_address = "fa:16:3e:f8:ab:da"
  }

  allowed_address_pairs {
    ip_address  = "10.0.0.202"
    mac_address = "fa:16:3e:ab:4b:58"
  }
}

data "vkcs_networking_port" "port_1" {
  name           = vkcs_networking_port.port_1.name
  admin_state_up = vkcs_networking_port.port_2.admin_state_up

  security_group_ids = [
    data.vkcs_networking_secgroup.default.id,
  ]
}

data "vkcs_networking_port" "port_2" {
  name           = vkcs_networking_port.port_1.name
  admin_state_up = vkcs_networking_port.port_2.admin_state_up

  tags = [
    "foo",
    "bar",
  ]
}

data "vkcs_networking_port" "port_3" {
  fixed_ip = vkcs_networking_port.port_1.all_fixed_ips.1
}
`
