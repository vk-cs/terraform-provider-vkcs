package vkcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSFSShareNetworkDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckSFS(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSShareNetworkDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.vkcs_sharedfilesystem_sharenetwork.sharenetwork_1", "id", "vkcs_sharedfilesystem_sharenetwork.sharenetwork_1", "id"),
					resource.TestCheckResourceAttr("data.vkcs_sharedfilesystem_sharenetwork.sharenetwork_1", "security_service_ids.#", "2"),
					resource.TestCheckResourceAttrPair("data.vkcs_sharedfilesystem_sharenetwork.sharenetwork_1", "neutron_net_id", "vkcs_sharedfilesystem_sharenetwork.sharenetwork_1", "neutron_net_id"),
					resource.TestCheckResourceAttrPair("data.vkcs_sharedfilesystem_sharenetwork.sharenetwork_1", "ip_version", "vkcs_sharedfilesystem_sharenetwork.sharenetwork_1", "ip_version"),
					resource.TestCheckResourceAttrPair("data.vkcs_sharedfilesystem_sharenetwork.sharenetwork_2", "id", "vkcs_sharedfilesystem_sharenetwork.sharenetwork_2", "id"),
					resource.TestCheckResourceAttr("data.vkcs_sharedfilesystem_sharenetwork.sharenetwork_2", "security_service_ids.#", "1"),
					resource.TestCheckResourceAttrPair("data.vkcs_sharedfilesystem_sharenetwork.sharenetwork_2", "neutron_net_id", "vkcs_sharedfilesystem_sharenetwork.sharenetwork_2", "neutron_net_id"),
					resource.TestCheckResourceAttrPair("data.vkcs_sharedfilesystem_sharenetwork.sharenetwork_2", "ip_version", "vkcs_sharedfilesystem_sharenetwork.sharenetwork_2", "ip_version"),
				),
			},
		},
	})
}

const testAccSFSShareNetworkDataSourceBasic = `
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

resource "vkcs_sharedfilesystem_securityservice" "securityservice_1" {
	name        = "security"
	description = "created by terraform"
	type        = "active_directory"
	server      = "192.168.199.10"
	dns_ip      = "192.168.199.10"
	domain      = "example.com"
	ou          = "CN=Computers,DC=example,DC=com"
	user        = "joinDomainUser"
	password    = "s8cret"
}

resource "vkcs_sharedfilesystem_securityservice" "securityservice_2" {
	name        = "security_through_obscurity"
	description = ""
	type        = "kerberos"
	server      = "192.168.199.11"
	dns_ip      = "192.168.199.11"
}

resource "vkcs_sharedfilesystem_sharenetwork" "sharenetwork_1" {
	name                = "test_sharenetwork_secure"
	description         = "share the secure love"
	neutron_net_id      = "${vkcs_networking_network.network_1.id}"
	neutron_subnet_id   = "${vkcs_networking_subnet.subnet_1.id}"
	security_service_ids = [
		"${vkcs_sharedfilesystem_securityservice.securityservice_1.id}",
		"${vkcs_sharedfilesystem_securityservice.securityservice_2.id}",
	]
}

resource "vkcs_sharedfilesystem_sharenetwork" "sharenetwork_2" {
	name                = "test_sharenetwork_secure"
	description         = "share the less secure love"
	neutron_net_id      = "${vkcs_networking_network.network_1.id}"
	neutron_subnet_id   = "${vkcs_networking_subnet.subnet_1.id}"
	security_service_ids = [
		"${vkcs_sharedfilesystem_securityservice.securityservice_1.id}",
	]
}

data "vkcs_sharedfilesystem_sharenetwork" "sharenetwork_1" {
	name                = "${vkcs_sharedfilesystem_sharenetwork.sharenetwork_1.name}"
	security_service_id = "${vkcs_sharedfilesystem_securityservice.securityservice_2.id}"
	ip_version          = "${vkcs_sharedfilesystem_sharenetwork.sharenetwork_1.ip_version}"
}

data "vkcs_sharedfilesystem_sharenetwork" "sharenetwork_2" {
	name                = "test_sharenetwork_secure"
	description         = "${vkcs_sharedfilesystem_sharenetwork.sharenetwork_2.description}"
	security_service_id = "${vkcs_sharedfilesystem_securityservice.securityservice_1.id}"
}
`
