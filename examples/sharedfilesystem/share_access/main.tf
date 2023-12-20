resource "vkcs_sharedfilesystem_share_access" "opencloud" {
  share_id     = vkcs_sharedfilesystem_share.data.id
  access_type  = "ip"
  access_to    = "192.168.199.11"
  access_level = "rw"
}
