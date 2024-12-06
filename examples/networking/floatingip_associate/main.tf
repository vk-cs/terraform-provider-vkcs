resource "vkcs_networking_floatingip_associate" "floatingip_associate" {
  floating_ip = vkcs_networking_floatingip.base_fip.address
  port_id     = vkcs_networking_port.persistent_etcd.id
  # Ensure the router interface is up
  depends_on = [vkcs_networking_router_interface.db]
}
