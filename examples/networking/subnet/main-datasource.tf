data "vkcs_networking_subnet" "subnet_one_of_internal" {
  cidr       = "192.168.199.0/24"
  network_id = vkcs_networking_network.app.id
  # This is unnecessary in real life.
  # This is required here to let the example work with subnet resource example. 
  depends_on = [vkcs_networking_subnet.app]
}
