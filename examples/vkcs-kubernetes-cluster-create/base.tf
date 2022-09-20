resource "vkcs_networking_network" "k8s" {
  name           = "k8s-net"
  admin_state_up = true
}

resource "vkcs_networking_subnet" "k8s-subnetwork" {
  name            = "k8s-subnet"
  network_id      = vkcs_networking_network.k8s.id
  cidr            = "192.168.199.0/24"
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}

data "vkcs_networking_network" "extnet" {
  name = "ext-net"
}

resource "vkcs_networking_router" "k8s" {
  name                = "k8s-router"
  admin_state_up      = true
  external_network_id = data.vkcs_networking_network.extnet.id
}

resource "vkcs_networking_router_interface" "k8s" {
  router_id = vkcs_networking_router.k8s.id
  subnet_id = vkcs_networking_subnet.k8s-subnetwork.id
}

data "vkcs_compute_flavor" "k8s" {
  name = "Standard-2-4-50"
}
