resource "vkcs_networking_router_route" "compute_gateway" {
  router_id        = vkcs_networking_router.router.id
  destination_cidr = "10.10.0.0/23"
  next_hop         = vkcs_compute_instance.gateway.access_ip_v4
}
