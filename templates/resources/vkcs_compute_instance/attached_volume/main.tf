resource "vkcs_blockstorage_volume" "myvol" {
  name = "myvol"
  size = 1
}

resource "vkcs_compute_instance" "myinstance" {
  name            = "myinstance"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "3"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  network {
    name = "my_network"
  }
}

resource "vkcs_compute_volume_attach" "attached" {
  instance_id = "${vkcs_compute_instance.myinstance.id}"
  volume_id   = "${vkcs_blockstorage_volume.myvol.id}"
}
