# Create networks
resource "vkcs_networking_network" "db" {
  name        = "db-tf-example"
  description = "Database network"
}

resource "vkcs_networking_subnet" "db" {
  name       = "db-tf-example"
  network_id = vkcs_networking_network.db.id
  cidr       = "192.168.166.0/24"
}

# Get external network with Internet access
data "vkcs_networking_network" "extnet" {
  name = "internet"
}

# Create a router to connect networks
resource "vkcs_networking_router" "router" {
  name = "router-tf-example"
  # Connect router to Internet
  external_network_id = data.vkcs_networking_network.extnet.id
}

# Connect networks to the router
resource "vkcs_networking_router_interface" "db" {
  router_id = vkcs_networking_router.router.id
  subnet_id = vkcs_networking_subnet.db.id
}

