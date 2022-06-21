resource "vkcs_networking_network" "network" {
  name           = "net"
}

resource "vkcs_networking_subnet" "subnetwork" {
  name       = "subnet_1"
  network_id = vkcs_networking_network.network.id
  cidr       = "192.168.199.0/24"
  ip_version = 4
}
