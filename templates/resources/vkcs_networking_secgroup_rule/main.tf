resource "vkcs_networking_secgroup" "secgroup_1" {
  name        = "secgroup_1"
  description = "My security group"
}

resource "vkcs_networking_secgroup_rule" "secgroup_rule_1" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${vkcs_networking_secgroup.secgroup_1.id}"
}
