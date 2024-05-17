data "vkcs_networking_secgroup" "default_secgroup" {
  name = "default"
  sdn  = "neutron"
}
