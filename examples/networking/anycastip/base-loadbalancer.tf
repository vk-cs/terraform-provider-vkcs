resource "vkcs_lb_loadbalancer" "app1" {
  name          = "app-tf-example-1"
  description   = "Loadbalancer for resources/datasources testing"
  vip_subnet_id = vkcs_networking_subnet.app.id
  tags          = ["app-front-tf-example"]
}

resource "vkcs_lb_loadbalancer" "app2" {
  name          = "app-tf-example-2"
  description   = "Loadbalancer for resources/datasources testing"
  vip_subnet_id = vkcs_networking_subnet.app.id
  tags          = ["app-front-tf-example"]
}


