resource "vkcs_networking_network" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "vkcs_compute_instance" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "vkcs_compute_interface_attach" "ai_1" {
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  network_id  = "${vkcs_networking_port.network_1.id}"
  fixed_ip    = "10.0.10.10"
}
