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
  security_group_id = "${vkcs_networking_secgroup.secgroup.id}"
	description = "secgroup_rule_1"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_2" {
  direction = "ingress"
  ethertype = "IPv4"
  port_range_max = 80
  port_range_min = 80
  protocol = "tcp"
  security_group_id = "${vkcs_networking_secgroup.secgroup.id}"
}
