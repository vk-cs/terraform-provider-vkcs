resource "vkcs_networking_port_secgroup_associate" "replace_secgroups" {
  port_id = vkcs_networking_port.persistent_etcd.id
  security_group_ids = [
    vkcs_networking_secgroup.http.id,
  ]
  enforce = true
}
