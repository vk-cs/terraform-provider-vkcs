resource "vkcs_lb_l7rule" "app_api_redirect" {
  l7policy_id  = vkcs_lb_l7policy.app_redirect.id
  compare_type = "EQUAL_TO"
  type         = "PATH"
  value        = "/api"
}
