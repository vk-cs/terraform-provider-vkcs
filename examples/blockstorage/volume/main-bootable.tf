resource "vkcs_blockstorage_volume" "bootable" {
  name              = "bootable-tf-example"
  size              = 5
  volume_type       = "ceph-ssd"
  image_id          = data.vkcs_images_image.debian.id
  availability_zone = "GZ1"
}
