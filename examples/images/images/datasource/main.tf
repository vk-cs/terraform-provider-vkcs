data "vkcs_images_images" "images" {
  visibility = "public"
  default = true
  properties = {
    mcs_os_distro = "debian"
  }
}
