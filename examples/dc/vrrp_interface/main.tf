resource "vkcs_dc_vrrp_interface" "dc_vrrp_interface" {
    name = "tf-example"
    description = "tf-example-description"
    dc_vrrp_id = vkcs_dc_vrrp.dc_vrrp.id
    dc_interface_id = vkcs_dc_interface.dc_interface.id
    priority = 100
    preempt = true
    master = true
}
