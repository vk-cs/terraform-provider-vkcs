resource "vkcs_blockstorage_volume" "volumes" {
  count             = 2
  name              = "new-volume-${count.index}-tf-example"
  size              = 1
  availability_zone = "GZ1"
  volume_type       = "ceph-ssd"
}
