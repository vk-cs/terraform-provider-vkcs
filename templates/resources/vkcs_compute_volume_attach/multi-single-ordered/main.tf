resource "vkcs_blockstorage_volume" "volumes" {
  count = 2
  name  = "${format("vol-%02d", count.index + 1)}"
  size  = 1
}

resource "vkcs_compute_instance" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "vkcs_compute_volume_attach" "attach_1" {
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  volume_id   = "${vkcs_blockstorage_volume.volumes.0.id}"
}

resource "vkcs_compute_volume_attach" "attach_2" {
  instance_id = "${vkcs_compute_instance.instance_1.id}"
  volume_id   = "${vkcs_blockstorage_volume.volumes.1.id}"

  depends_on = ["vkcs_compute_volume_attach.attach_1"]
}

output "volume_devices" {
  value = "${vkcs_compute_volume_attach.attachments.*.device}"
}
