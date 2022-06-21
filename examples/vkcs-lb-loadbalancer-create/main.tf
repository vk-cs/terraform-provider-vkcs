resource "vkcs_compute_instance" "compute_1" {
  name            = "compute-instance-1"
  flavor_id       = data.vkcs_compute_flavor.compute.id
  security_groups = ["default"]
  image_id = data.vkcs_images_image.compute.id

  network {
    uuid = vkcs_networking_network.lb.id
    fixed_ip_v4 = "192.168.199.110"
  }

  depends_on = [
    vkcs_networking_network.lb,
    vkcs_networking_subnet.lb
  ]
}

resource "vkcs_compute_instance" "compute_2" {
  name            = "compute-instance-2"
  flavor_id       = data.vkcs_compute_flavor.compute.id
  security_groups = ["default"]
  image_id = data.vkcs_images_image.compute.id

  network {
    uuid = vkcs_networking_network.lb.id
    fixed_ip_v4 = "192.168.199.111"
  }

  depends_on = [
    vkcs_networking_network.lb,
    vkcs_networking_subnet.lb
  ]
}

resource "vkcs_lb_loadbalancer" "loadbalancer" {
  name = "loadbalancer"
  vip_subnet_id = "${vkcs_networking_subnet.lb.id}"
  tags = ["tag1"]
}

resource "vkcs_lb_listener" "listener" {
  name = "listener"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${vkcs_lb_loadbalancer.loadbalancer.id}"
}

resource "vkcs_lb_pool" "pool" {
  name = "pool"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  listener_id = "${vkcs_lb_listener.listener.id}"
}

resource "vkcs_lb_member" "member_1" {
  address = "192.168.199.110"
  protocol_port = 8080
  pool_id = "${vkcs_lb_pool.pool.id}"
  subnet_id = "${vkcs_networking_subnet.lb.id}"
  weight = 0
}

resource "vkcs_lb_member" "member_2" {
  address = "192.168.199.111"
  protocol_port = 8080
  pool_id = "${vkcs_lb_pool.pool.id}"
  subnet_id = "${vkcs_networking_subnet.lb.id}"
}
