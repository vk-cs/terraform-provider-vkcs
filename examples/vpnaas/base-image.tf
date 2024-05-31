data "vkcs_images_image" "image" {
  properties = {
    mcs_os_distro  = "centos"
    mcs_os_version = "7.9"
  }
  default    = true
  visibility = "public"
}
