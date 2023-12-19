resource "vkcs_compute_volume_attach" "data" {
  instance_id = vkcs_compute_instance.basic.id
  volume_id   = vkcs_blockstorage_volume.data.id
}
