resource "vkcs_networking_floatingip" "myip" {
  pool = "my_pool"
}

resource "vkcs_compute_instance" "multi-net" {
  name            = "multi-net"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  network {
    name = "my_first_network"
  }

  network {
    name = "my_second_network"
  }
}

resource "vkcs_compute_floatingip_associate" "myip" {
  floating_ip = "${vkcs_networking_floatingip.myip.address}"
  instance_id = "${vkcs_compute_instance.multi-net.id}"
  fixed_ip    = "${vkcs_compute_instance.multi-net.network.1.fixed_ip_v4}"
}
