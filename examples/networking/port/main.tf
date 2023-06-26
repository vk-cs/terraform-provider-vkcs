resource "vkcs_networking_port" "persistent_etcd" {
  name       = "persistent-etcd-tf-example"
  network_id = vkcs_networking_network.db.id
  # Specify subnet for multi subnet networks to controls
  # which subnet is used to get port IP
  # Also this brings you required dependency of the port
  # Otherwise if you create the port with network and subnet
  # in one tf file the port may be created before the subnet
  # So it does not get IP
  # Alternative for this case is to set subnet dependency
  # explicitly
  fixed_ip {
    subnet_id = vkcs_networking_subnet.db.id
  }
  # Specify required security groups instead of getting 'default' one
  security_group_ids = [vkcs_networking_secgroup.etcd.id]
}
