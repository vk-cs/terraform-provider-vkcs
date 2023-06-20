resource "vkcs_blockstorage_volume" "volume" {
  name = "volume"
  description = "test volume"
  metadata = {
    foo = "bar"
  }
  size = 1
  availability_zone = "GZ1"
  volume_type = "ceph-ssd"
}

resource "vkcs_blockstorage_snapshot" "snapshot" {
  volume_id = "${vkcs_blockstorage_volume.volume.id}"
  name = "snapshot"
  description = "test snapshot"
  metadata = {
    foo = "bar"
  }
}
