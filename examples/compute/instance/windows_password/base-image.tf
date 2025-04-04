data "vkcs_images_image" "windows" {
  visibility  = "public"
  default     = true
  most_recent = true
  properties = {
    mcs_os_type = "windows"
    os_version  = "10.0"
  }
}
