data "vkcs_images_image" "ubuntu" {
  name        = "Ubuntu 16.04"
  most_recent = true

  properties = {
    key = "value"
  }
}
