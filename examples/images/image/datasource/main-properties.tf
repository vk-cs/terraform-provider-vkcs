data "vkcs_images_image" "ubuntu" {
  visibility = "public"
  default  = true
  properties = {
    mcs_os_distro  = "ubuntu"
    mcs_os_version = "22.04"
  }
}
