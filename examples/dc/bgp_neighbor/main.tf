resource "vkcs_dc_bgp_neighbor" "dc_bgp_neighbor" {
    name = "tf-example"
    add_paths = "on"
    description = "tf-example-description"
    dc_bgp_id = vkcs_dc_bgp_instance.dc_bgp_instance.id
    remote_asn = 1
    remote_ip = "192.168.1.3"
    enabled = true
}
