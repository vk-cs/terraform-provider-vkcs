data "vkcs_images_image" "eurolinux9" {
  tag         = "tf-example"
  # Useful if you keep old versions of your images.
  most_recent = true
  properties = {
    mcs_os_distro  = "eurolinux"
    mcs_os_version = "9"
  }
  # This is unnecessary in real life.
  # This is required here to let the example work with image resource example.
  depends_on = [ vkcs_images_image.eurolinux9 ]
}
