resource "vkcs_blockstorage_snapshot" "recent_snapshot" {
  volume_id = vkcs_blockstorage_volume.data.id
  name = "snapshot-tf-example"
  description = "test snapshot"
  metadata = {
    foo = "bar"
  }
}
