resource "vkcs_vpnaas_service" "vpn_to_datacenter" {
  name      = "vpn-tf-example"

  # See the argument description and check vkcs_networks_sdn datasource output to figure out
  # what type of router you should use in certain case (router or dc_router)
  router_id = vkcs_networking_router.router.id
}
