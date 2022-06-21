data "vkcs_compute_flavor" "compute" {
  name = "Basic-1-2-20"
}

data "vkcs_networking_network" "extnet" {
  name = "ext-net"
}

data "vkcs_images_image" "compute" {
  name = "Ubuntu-18.04-Standard"
}

resource "vkcs_networking_network" "compute" {
  name = "compute-net"
}

resource "vkcs_networking_subnet" "compute" {
  name       = "subnet_1"
  network_id = vkcs_networking_network.compute.id
  cidr       = "192.168.199.0/24"
  ip_version = 4
}

resource "vkcs_networking_router" "compute" {
  name                = "db-router"
  admin_state_up      = true
  external_network_id = data.vkcs_networking_network.extnet.id
}

resource "vkcs_networking_router_interface" "compute" {
  router_id = vkcs_networking_router.compute.id
  subnet_id = vkcs_networking_subnet.compute.id
}
