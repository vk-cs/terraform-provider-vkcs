data "vkcs_networking_port" "persistent_etcd" {
  tags       = ["tf-example", "etcd"]
  network_id = vkcs_networking_network.db.id
  # This is unnecessary in real life.
  # This is required here to let the example work with port resource example. 
  depends_on = [vkcs_networking_port.persistent_etcd]
}
