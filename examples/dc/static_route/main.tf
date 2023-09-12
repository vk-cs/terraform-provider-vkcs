resource "vkcs_dc_static_route" "dc_static_route" {
    name = "tf-example"
    description = "tf-example-description"
    dc_router_id = vkcs_dc_router.dc_router.id
    network = "192.168.1.0/24"
    gateway = "192.168.1.3"
    metric = 1
}
