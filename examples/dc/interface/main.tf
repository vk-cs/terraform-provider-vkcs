# Connect networks to the router
resource "vkcs_dc_interface" "dc_interface" {
  name         = "tf-example"
  description  = "tf-example-description"
  dc_router_id = vkcs_dc_router.dc_router.id
  network_id   = vkcs_networking_network.app.id
  subnet_id    = vkcs_networking_subnet.app.id
}
