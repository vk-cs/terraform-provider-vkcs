resource "vkcs_networking_secgroup" "ssh" {
  name = "ssh-tf-example"
}

resource "vkcs_networking_secgroup_rule" "ssh" {
  description       = "SSH rule"
  security_group_id = vkcs_networking_secgroup.ssh.id
  direction         = "ingress"
  protocol          = "tcp"
  # Specify SSH port
  port_range_max = 22
  port_range_min = 22
  # Allow access from any sources
  remote_ip_prefix = "0.0.0.0/0"
}

data "vkcs_networking_secgroup" "default" {
  name = "default"
  sdn  = "sprut"
}

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
