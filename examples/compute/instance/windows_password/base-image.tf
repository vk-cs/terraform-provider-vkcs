data "vkcs_images_image" "windows" {
  visibility  = "public"
  default     = true
  most_recent = true
  properties = {
    mcs_os_type    = "windows"
    mcs_os_version = "2022"
    mcs_os_edition = "std"
    mcs_os_lang    = "en"
  }
}
