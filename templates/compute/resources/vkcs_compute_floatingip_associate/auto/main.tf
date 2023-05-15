resource "vkcs_compute_instance" "instance_1" {
  name            = "instance_1"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = 3
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]
}

resource "vkcs_networking_floatingip" "fip_1" {
  pool = "my_pool"
}

resource "vkcs_compute_floatingip_associate" "fip_1" {
  floating_ip = "${vkcs_networking_floatingip.fip_1.address}"
  instance_id = "${vkcs_compute_instance.instance_1.id}"
}
