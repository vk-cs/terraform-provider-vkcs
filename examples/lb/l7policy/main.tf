resource "vkcs_lb_l7policy" "app_redirect" {
  name             = "http-tf-example"
  description      = "Policy for tf lb testing"
  action           = "REDIRECT_TO_POOL"
  position         = 1
  listener_id      = vkcs_lb_listener.app_http.id
  redirect_pool_id = vkcs_lb_pool.http.id
}
