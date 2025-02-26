# Create networks
resource "vkcs_networking_network" "app" {
  name        = "app-tf-example"
  description = "Application network"
  sdn         = "neutron"
}

resource "vkcs_networking_subnet" "app" {
  name       = "app-tf-example"
  network_id = vkcs_networking_network.app.id
  cidr       = "192.168.199.0/24"
  sdn        = "neutron"
}

# Get external network with Internet access
data "vkcs_networking_network" "extnet" {
  name = "ext-net"
  sdn  = "neutron"
}
