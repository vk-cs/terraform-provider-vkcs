resource "vkcs_blockstorage_volume" "myvol" {
  name     = "myvol"
  size     = 5
  image_id = "<image-id>"
}

resource "vkcs_compute_instance" "boot-from-volume" {
  name            = "bootfromvolume"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    uuid                  = "${vkcs_blockstorage_volume.myvol.id}"
    source_type           = "volume"
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
  }

  network {
    name = "my_network"
  }
}
