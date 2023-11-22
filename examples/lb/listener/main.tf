resource "vkcs_lb_listener" "app_http" {
  name            = "app-http-tf-example"
  description     = "Listener for resources/datasources testing"
  loadbalancer_id = vkcs_lb_loadbalancer.app.id
  protocol        = "HTTP"
  protocol_port   = 8080
}
