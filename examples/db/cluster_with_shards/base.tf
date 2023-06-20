data "vkcs_compute_flavor" "db" {
  name = "Standard-2-8-50"
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

resource "vkcs_networking_secgroup" "secgroup" {
  name = "security_group"
  description = "terraform security group"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_1" {
  direction = "ingress"
  ethertype = "IPv4"
  port_range_max = 22
  port_range_min = 22
  protocol = "tcp"
  remote_ip_prefix = "0.0.0.0/0"
  security_group_id = vkcs_networking_secgroup.secgroup.id
	description = "secgroup_rule_1"
}
