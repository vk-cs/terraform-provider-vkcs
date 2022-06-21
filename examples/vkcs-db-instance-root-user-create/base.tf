data "vkcs_compute_flavor" "db" {
  name = "Basic-1-2-20"
}

data "vkcs_networking_network" "extnet" {
  name = "ext-net"
}

resource "vkcs_networking_network" "db" {
  name           = "db-net"
  admin_state_up = true
}

resource "vkcs_networking_subnet" "db" {
  name       = "subnet_1"
  network_id = vkcs_networking_network.db.id
  cidr       = "192.168.199.0/24"
  ip_version = 4
}

resource "vkcs_networking_router" "db" {
  name                = "db-router"
  admin_state_up      = true
  external_network_id = data.vkcs_networking_network.extnet.id
}

resource "vkcs_networking_router_interface" "db" {
  router_id = vkcs_networking_router.db.id
  subnet_id = vkcs_networking_subnet.db.id
}
