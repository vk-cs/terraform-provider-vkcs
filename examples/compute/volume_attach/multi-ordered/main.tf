resource "vkcs_compute_volume_attach" "attach_1" {
  instance_id = vkcs_compute_instance.basic.id
  volume_id   = vkcs_blockstorage_volume.volumes.0.id
}

resource "vkcs_compute_volume_attach" "attach_2" {
  instance_id = vkcs_compute_instance.basic.id
  volume_id   = vkcs_blockstorage_volume.volumes.1.id

  depends_on = [vkcs_compute_volume_attach.attach_1]
}
