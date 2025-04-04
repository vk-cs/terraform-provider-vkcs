resource "vkcs_networking_secgroup" "rdp" {
  name = "rdp-tf-example"
}

resource "vkcs_networking_secgroup_rule" "rdp" {
  description       = "RDP rule"
  security_group_id = vkcs_networking_secgroup.rdp.id
  direction         = "ingress"
  protocol          = "tcp"
  # Specify RDP port
  port_range_max = 3389
  port_range_min = 3389
  # Allow access from any sources
  remote_ip_prefix = "0.0.0.0/0"
}
