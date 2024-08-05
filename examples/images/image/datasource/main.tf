data "vkcs_images_image" "debian" {
  # Both arguments are required to search an actual image provided by VKCS.
  visibility = "public"
  default    = true
  # Use properties to distinguish between available images.
  properties = {
    mcs_os_distro  = "debian"
    mcs_os_version = "12"
  }
}
