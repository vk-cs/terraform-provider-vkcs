resource "vkcs_dc_bgp_instance" "dc_bgp_instance" {
    name = "tf-example"
    description = "tf-example-description"
    dc_router_id = vkcs_dc_router.dc_router.id
    bgp_router_id = "192.168.1.2"
    asn = 12345
    ecmp_enabled = true
    enabled = true
    graceful_restart = true
}
