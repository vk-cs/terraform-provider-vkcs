resource "vkcs_dc_vrrp_address" "dc_vrrp_address" {
  name        = "tf-example"
  description = "tf-example-description"
  dc_vrrp_id  = vkcs_dc_vrrp.dc_vrrp.id
  ip_address  = "192.168.199.42"
}
