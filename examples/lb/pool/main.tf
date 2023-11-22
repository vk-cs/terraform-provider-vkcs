resource "vkcs_lb_pool" "http" {
  name        = "http-tf-example"
  description = "Pool for http member/members testing"
  listener_id = vkcs_lb_listener.app_http.id
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
}
