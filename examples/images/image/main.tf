resource "vkcs_images_image" "eurolinux9" {
  name             = "eurolinux9-tf-example"
  image_source_url = "https://fbi.cdn.euro-linux.com/images/EL-9-cloudgeneric-2023-03-19.raw.xz"
  # compression_format should be set for compressed image source.
  compression_format = "xz"
  container_format   = "bare"
  disk_format        = "raw"
  # Minimal requirements from image vendor.
  # Should be set to prevent VKCS to build VM on lesser resources.
  min_ram_mb  = 1536
  min_disk_gb = 10

  properties = {
    # Refer to https://mcs.mail.ru/docs/en/base/iaas/instructions/vm-images/vm-image-metadata
    os_type             = "linux"
    os_admin_user       = "root"
    mcs_name            = "EuroLinux 9"
    mcs_os_distro       = "eurolinux"
    mcs_os_version      = "9"
    hw_qemu_guest_agent = "yes"
    os_require_quiesce  = "yes"
  }
  # Use tags to organize your images.
  tags = ["tf-example"]
}
