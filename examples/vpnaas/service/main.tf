resource "vkcs_vpnaas_service" "vpn_to_datacenter" {
  name      = "vpn-tf-example"
  router_id = vkcs_networking_router.router.id
}
