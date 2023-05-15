resource "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "vkcs_compute_instance" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "vkcs_compute_volume_attach" "va_1" {
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  volume_id   = "${vkcs_blockstorage_volume.volume_1.id}"
}
