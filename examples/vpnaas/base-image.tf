data "vkcs_images_image" "image" {
  visibility = "public"
  default    = true
  properties = {
    mcs_os_distro  = "debian"
    mcs_os_version = "12"
  }
}
