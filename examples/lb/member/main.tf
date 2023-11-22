resource "vkcs_lb_member" "front_http" {
  pool_id       = vkcs_lb_pool.http.id
  address       = "192.168.199.110"
  protocol_port = 8080
}
