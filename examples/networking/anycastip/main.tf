resource "vkcs_networking_anycastip" "anycastip" {
  name        = "app-tf-example"
  description = "app-tf-example"

  network_id = data.vkcs_networking_network.extnet.id
  associations = [
    {
      id   = vkcs_lb_loadbalancer.app1.vip_port_id
      type = "octavia"
    },
    {
      id   = vkcs_lb_loadbalancer.app2.vip_port_id
      type = "octavia"
    }
  ]
}
