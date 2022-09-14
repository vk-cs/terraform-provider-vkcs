resource "vkcs_compute_flavor" "flavor-with-swap" {
  name  = "flavor-with-swap"
  ram   = "8096"
  vcpus = "2"
  disk  = "20"
  swap  = "4096"
}
resource "vkcs_compute_instance" "vm-swap" {
  name            = "vm_swap"
  flavor_id       = "${vkcs_compute_flavor.flavor-with-swap.id}"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]
  block_device {
    boot_index            = 0
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "image"
    uuid                  = "<image-id>"
  }
  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "blank"
    guest_format          = "swap"
    volume_size           = 4
  }
  network {
    name = "my_network"
  }
}
