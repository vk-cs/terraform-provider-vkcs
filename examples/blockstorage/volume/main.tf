resource "vkcs_blockstorage_volume" "data" {
  name = "data-tf-example"
  description = "test volume"
  metadata = {
    foo = "bar"
  }
  size = 1
  availability_zone = "GZ1"
  volume_type = "ceph-ssd"
}
