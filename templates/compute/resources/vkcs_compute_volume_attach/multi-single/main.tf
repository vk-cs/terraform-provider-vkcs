resource "vkcs_blockstorage_volume" "volumes" {
  count = 2
  name  = "${format("vol-%02d", count.index + 1)}"
  size  = 1
}

resource "vkcs_compute_instance" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "vkcs_compute_volume_attach" "attachments" {
  count       = 2
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  volume_id   = "${vkcs_blockstorage_volume.volumes.*.id[count.index]}"
}

output "volume_devices" {
  value = "${vkcs_compute_volume_attach.attachments.*.device}"
}
