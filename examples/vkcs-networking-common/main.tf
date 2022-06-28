data "vkcs_networking_network" "extnet" {
  name = "ext-net"
}

resource "vkcs_networking_network" "network" {
  name = "net"
}

resource "vkcs_networking_subnet" "subnetwork" {
  name       = "subnet_1"
  network_id = vkcs_networking_network.network.id
  cidr       = "192.168.199.0/24"
  ip_version = 4
}

resource "vkcs_networking_router" "router" {
  name                = "router"
  admin_state_up      = true
  external_network_id = data.vkcs_networking_network.extnet.id
}

resource "vkcs_networking_router_interface" "db" {
  router_id = vkcs_networking_router.router.id
  subnet_id = vkcs_networking_subnet.subnetwork.id
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

resource "vkcs_networking_secgroup_rule" "secgroup_rule_2" {
  direction = "ingress"
  ethertype = "IPv4"
  port_range_max = 3389
  port_range_min = 3389
  remote_ip_prefix = "0.0.0.0/0"
  protocol = "tcp"
  security_group_id = vkcs_networking_secgroup.secgroup.id
}

resource "vkcs_networking_port" "port" {
  name = "port_1"
  admin_state_up = "true"
  network_id = vkcs_networking_network.network.id

  fixed_ip {
    subnet_id =  vkcs_networking_subnet.subnetwork.id
    ip_address = "192.168.199.23"
  }
}

resource "vkcs_networking_port_secgroup_associate" "port" {
  port_id = vkcs_networking_port.port.id
  enforce = "false"
  security_group_ids = [
    vkcs_networking_secgroup.secgroup.id,
  ]
}
