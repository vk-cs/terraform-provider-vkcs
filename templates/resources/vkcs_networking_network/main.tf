resource "vkcs_networking_network" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name       = "subnet_1"
  network_id = "${vkcs_networking_network.network_1.id}"
  cidr       = "192.168.199.0/24"
  ip_version = 4
}

resource "vkcs_compute_secgroup" "secgroup_1" {
  name        = "secgroup_1"
  description = "a security group"

  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
}

resource "vkcs_networking_port" "port_1" {
  name               = "port_1"
  network_id         = "${vkcs_networking_network.network_1.id}"
  admin_state_up     = "true"
  security_group_ids = ["${vkcs_compute_secgroup.secgroup_1.id}"]

  fixed_ip {
    "subnet_id"  = "${vkcs_networking_subnet.subnet_1.id}"
    "ip_address" = "192.168.199.10"
  }
}

resource "vkcs_compute_instance" "instance_1" {
  name            = "instance_1"
  security_groups = ["${vkcs_compute_secgroup.secgroup_1.name}"]

  network {
    port = "${vkcs_networking_port.port_1.id}"
  }
}
