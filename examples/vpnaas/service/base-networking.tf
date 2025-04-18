# Create sprut network
resource "vkcs_networking_network" "sprut" {
  name = "sprut-network-tf-example"
}

resource "vkcs_networking_subnet" "sprut" {
  name       = "sprut-subnet-tf-example"
  network_id = vkcs_networking_network.sprut.id
  cidr       = "172.16.0.0/24"
}

# Get external network with Internet access
data "vkcs_networking_network" "internet_sprut" {
  name = "internet"
}

# Create a router to connect networks
resource "vkcs_dc_router" "router" {
  availability_zone = "GZ1"
  flavor            = "standard"
  name              = "dc-router-sprut-tf-example"
  description       = "dc_router in sprut"
}

# Connect the router to Internet
resource "vkcs_dc_interface" "internet_sprut" {
  name         = "dc-interface-for-internet-sprut-tf-example"
  description  = "dc_interface for connecting dc_router to the internet"
  dc_router_id = vkcs_dc_router.router.id
  network_id   = data.vkcs_networking_network.internet_sprut.id
}

# Connect networks to the router
resource "vkcs_dc_interface" "subnet_sprut" {
  name         = "dc-interface-for-subnet-sprut-tf-example"
  description  = "dc_interface for connecting dc_router to the network and subnet"
  dc_router_id = vkcs_dc_router.router.id
  network_id   = vkcs_networking_network.sprut.id
  subnet_id    = vkcs_networking_subnet.sprut.id
}
