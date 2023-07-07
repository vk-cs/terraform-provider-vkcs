data "vkcs_images_image" "compute" {
  name = "Ubuntu-18.04-Standard"
}

data "vkcs_compute_flavor" "compute" {
  name = "Basic-1-2-20"
}

resource "vkcs_networking_network" "lb" {
  name = "network"
}

resource "vkcs_networking_subnet" "lb" {
  name       = "subnet"
  cidr       = "192.168.199.0/24"
  network_id = vkcs_networking_network.lb.id
}
