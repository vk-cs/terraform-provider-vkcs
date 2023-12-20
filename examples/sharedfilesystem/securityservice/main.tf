resource "vkcs_sharedfilesystem_securityservice" "ad_common" {
  name        = "active-directory-tf-example"
  description = "active directory tf example"
  type        = "active_directory"
  server      = "192.168.199.10"
  dns_ip      = "192.168.199.10"
}
