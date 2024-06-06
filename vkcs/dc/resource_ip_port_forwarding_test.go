package dc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/acctest"
)

func TestAccDCIPPortForwarding_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDCIPPortForwardingBasic, map[string]string{
					"TestAccDCInterfaceBasic": acctest.AccTestRenderConfig(testAccDCInterfaceBasic, map[string]string{
						"TestAccDCRouterBasic": acctest.AccTestRenderConfig(testAccDCRouterBasic),
					}),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_dc_ip_port_forwarding.dc_ip_port_forwarding", "protocol", "udp"),
					resource.TestCheckResourceAttr("vkcs_dc_ip_port_forwarding.dc_ip_port_forwarding", "to_destination", "172.17.20.30"),
				),
			},
			{
				ResourceName:      "vkcs_dc_ip_port_forwarding.dc_ip_port_forwarding",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDCIPPortForwarding_full(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV6ProviderFactories: acctest.AccTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.AccTestRenderConfig(testAccDCIPPortForwardingFull, map[string]string{
					"TestAccDCInterfaceBasic": acctest.AccTestRenderConfig(testAccDCInterfaceBasic, map[string]string{
						"TestAccDCRouterBasic": acctest.AccTestRenderConfig(testAccDCRouterBasic),
					}),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_dc_ip_port_forwarding.dc_ip_port_forwarding", "name", "tfacc-dc-ip-port-forwarding"),
					resource.TestCheckResourceAttr("vkcs_dc_ip_port_forwarding.dc_ip_port_forwarding", "description", "tfacc-dc-ip-port-forwarding-description"),
					resource.TestCheckResourceAttr("vkcs_dc_ip_port_forwarding.dc_ip_port_forwarding", "protocol", "udp"),
					resource.TestCheckResourceAttr("vkcs_dc_ip_port_forwarding.dc_ip_port_forwarding", "source", "192.168.124.0/32"),
					resource.TestCheckResourceAttr("vkcs_dc_ip_port_forwarding.dc_ip_port_forwarding", "destination", "10.10.0.100"),
					resource.TestCheckResourceAttr("vkcs_dc_ip_port_forwarding.dc_ip_port_forwarding", "port", "80"),
					resource.TestCheckResourceAttr("vkcs_dc_ip_port_forwarding.dc_ip_port_forwarding", "to_destination", "172.17.20.30"),
					resource.TestCheckResourceAttr("vkcs_dc_ip_port_forwarding.dc_ip_port_forwarding", "to_port", "80"),
				),
			},
			{
				Config: acctest.AccTestRenderConfig(testAccDCIPPortForwardingFullUpdate, map[string]string{
					"TestAccDCInterfaceBasic": acctest.AccTestRenderConfig(testAccDCInterfaceBasic, map[string]string{
						"TestAccDCRouterBasic": acctest.AccTestRenderConfig(testAccDCRouterBasic),
					}),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("vkcs_dc_ip_port_forwarding.dc_ip_port_forwarding", "name", "tfacc-dc-ip-port-forwarding-upd"),
					resource.TestCheckResourceAttr("vkcs_dc_ip_port_forwarding.dc_ip_port_forwarding", "description", "tfacc-dc-ip-port-forwarding-description-upd"),
					resource.TestCheckResourceAttr("vkcs_dc_ip_port_forwarding.dc_ip_port_forwarding", "protocol", "tcp"),
					resource.TestCheckResourceAttr("vkcs_dc_ip_port_forwarding.dc_ip_port_forwarding", "source", "192.168.125.0/32"),
					resource.TestCheckResourceAttr("vkcs_dc_ip_port_forwarding.dc_ip_port_forwarding", "destination", "10.10.0.101"),
					resource.TestCheckResourceAttr("vkcs_dc_ip_port_forwarding.dc_ip_port_forwarding", "port", "90"),
					resource.TestCheckResourceAttr("vkcs_dc_ip_port_forwarding.dc_ip_port_forwarding", "to_destination", "172.17.20.31"),
					resource.TestCheckResourceAttr("vkcs_dc_ip_port_forwarding.dc_ip_port_forwarding", "to_port", "90"),
				),
			},
		},
	})
}

const testAccDCIPPortForwardingBasic = `
{{ .TestAccDCInterfaceBasic}}

resource "vkcs_dc_ip_port_forwarding" "dc_ip_port_forwarding" {
	dc_interface_id = vkcs_dc_interface.dc_interface.id
	protocol = "udp"
    to_destination = "172.17.20.30"
}
`

const testAccDCIPPortForwardingFull = `
{{ .TestAccDCInterfaceBasic}}

resource "vkcs_dc_ip_port_forwarding" "dc_ip_port_forwarding" {
	name = "tfacc-dc-ip-port-forwarding"
	description = "tfacc-dc-ip-port-forwarding-description"
	dc_interface_id = vkcs_dc_interface.dc_interface.id
	protocol = "udp"
	source = "192.168.124.0/32"
	destination = "10.10.0.100"
	port = 80
    to_destination = "172.17.20.30"
	to_port = 80
}
`

const testAccDCIPPortForwardingFullUpdate = `
{{ .TestAccDCInterfaceBasic}}

resource "vkcs_dc_ip_port_forwarding" "dc_ip_port_forwarding" {
	name = "tfacc-dc-ip-port-forwarding-upd"
	description = "tfacc-dc-ip-port-forwarding-description-upd"
	dc_interface_id = vkcs_dc_interface.dc_interface.id
	protocol = "tcp"
	source = "192.168.125.0/32"
	destination = "10.10.0.101"
	port = 90
    to_destination = "172.17.20.31"
	to_port = 90
}
`
