data "vkcs_networking_network" "ext_network" {
  name = "public"
}

resource "vkcs_networking_floatingip" "floatip_1" {
  pool       = data.vkcs_networking_network.ext_network.name
  subnet_ids = [<subnet1_id>, <subnet2_id>, <subnet3_id>]
}