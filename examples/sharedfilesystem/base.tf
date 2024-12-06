resource "vkcs_networking_network" "sfs" {
  name = "network"
  sdn  = "sprut"
}

resource "vkcs_networking_subnet" "sfs" {
  name       = "subnet"
  cidr       = "192.168.199.0/24"
  network_id = vkcs_networking_network.sfs.id
}
