resource "vkcs_cloud_monitoring" "basic" {
  image_id = data.vkcs_images_image.debian.id
}
