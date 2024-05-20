# Connect internet to the router
resource "vkcs_dc_interface" "dc_interface_internet" {
  name         = "interface-for-internet"
  dc_router_id = vkcs_dc_router.dc_router.id
  network_id   = data.vkcs_networking_network.internet_sprut.id
}
