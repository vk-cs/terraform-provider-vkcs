resource "vkcs_dc_bgp_static_announce" "dc_bgp_static_announce" {
  name        = "tf-example"
  description = "tf-example-description"
  dc_bgp_id   = vkcs_dc_bgp_instance.dc_bgp_instance.id
  network     = "192.168.1.0/24"
  gateway     = "192.168.1.3"
}
