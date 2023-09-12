resource "vkcs_dc_vrrp" "dc_vrrp" {
    name = "tf-example"
    description = "tf-example-description"
    group_id = 100
    network_id = vkcs_networking_network.app.id
    subnet_id = vkcs_networking_subnet.app.id
    advert_interval = 1
    enabled = true
}
