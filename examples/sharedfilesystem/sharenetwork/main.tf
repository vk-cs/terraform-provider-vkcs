resource "vkcs_sharedfilesystem_sharenetwork" "data" {
  name              = "sharenetwork-tf-example"
  description       = "sharing network for tf example"
  neutron_net_id    = vkcs_networking_network.app.id
  neutron_subnet_id = vkcs_networking_subnet.app.id
}
