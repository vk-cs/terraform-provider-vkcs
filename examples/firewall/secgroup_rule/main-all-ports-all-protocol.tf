resource "vkcs_networking_secgroup_rule" "all" {
  description       = "Any inbound traffic from etcd hosts"
  security_group_id = vkcs_networking_secgroup.etcd.id
  direction         = "ingress"
  remote_group_id   = vkcs_networking_secgroup.etcd.id
}
