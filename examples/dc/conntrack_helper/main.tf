resource "vkcs_dc_conntrack_helper" "dc-conntrack-helper" {
  dc_router_id = vkcs_dc_router.dc_router.id
  name         = "tf-example"
  description  = "tf-example-description"
  helper       = "ftp"
  protocol     = "tcp"
  port         = 21
}
