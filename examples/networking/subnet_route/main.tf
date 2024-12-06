resource "vkcs_networking_subnet_route" "subnet_route_to_external_tf_example" {
  subnet_id        = vkcs_networking_subnet.app.id
  destination_cidr = "10.0.1.0/24"
  next_hop         = vkcs_networking_port.persistent_etcd.all_fixed_ips[0]
}
