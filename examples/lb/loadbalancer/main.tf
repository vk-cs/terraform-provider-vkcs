resource "vkcs_lb_loadbalancer" "app" {
  name          = "app-tf-example"
  description   = "Loadbalancer for resources/datasources testing"
  vip_subnet_id = vkcs_networking_subnet.app.id
  tags          = ["app-front-tf-example"]
}
