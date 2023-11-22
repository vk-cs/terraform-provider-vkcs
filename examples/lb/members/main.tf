resource "vkcs_lb_members" "front_workers" {
  pool_id = vkcs_lb_pool.http.id

  dynamic "member" {
    for_each = vkcs_compute_instance.front_worker
    content {
      address       = member.value.access_ip_v4
      protocol_port = 8080
    }
  }
}
