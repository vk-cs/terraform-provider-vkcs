resource "vkcs_networking_secgroup_rule" "all_udp" {
  description       = "All inbound UDP traffic from etcd hosts"
  security_group_id = vkcs_networking_secgroup.etcd.id
  direction         = "ingress"
  protocol          = "udp"
  remote_group_id   = vkcs_networking_secgroup.etcd.id
}
