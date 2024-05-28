data "vkcs_images_image" "image" {
  size_max    = 3 * 1024 * 1024
  visibility  = "public"
  most_recent = true
  properties = {
    mcs_os_distro  = "centos"
    mcs_os_version = "7.9"
  }
}
