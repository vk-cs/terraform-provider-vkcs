# Create networks
resource "vkcs_networking_network" "k8s_network" {
  name        = "k8s-network"
  sdn         = "sprut"
  description = "Network for service ManagedK8S."
}

resource "vkcs_networking_subnet" "k8s_subnet" {
  name       = "k8s-subnet"
  network_id = vkcs_networking_network.k8s_network.id
  cidr       = "192.168.199.0/24"
}

# Get external network with Internet access
data "vkcs_networking_network" "extnet" {
  name = "internet"
}

# Create a router to connect networks
resource "vkcs_networking_router" "k8s_router" {
  name = "k8s-router"
  # Connect router to the Internet
  external_network_id = data.vkcs_networking_network.extnet.id
}

# Connect subnet to the router
resource "vkcs_networking_router_interface" "extnet_conn" {
  router_id = vkcs_networking_router.k8s_router.id
  subnet_id = vkcs_networking_subnet.k8s_subnet.id
}
