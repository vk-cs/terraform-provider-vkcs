resource "vkcs_networking_port_secgroup_associate" "remove_secgroups" {
  port_id            = vkcs_networking_port.persistent_etcd.id
  security_group_ids = []
  enforce            = true
}
