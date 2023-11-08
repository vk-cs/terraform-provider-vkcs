data "vkcs_networking_floatingip" "fip_by_port" {
  port_id = vkcs_networking_port.persistent_etcd.id
  # This is unnecessary in real life.
  # This is required here to let the example work with floating ip resource example. 
  depends_on = [vkcs_networking_floatingip.associated_fip]
}
