resource "vkcs_sharedfilesystem_share" "data" {
  name             = "share-data-tf-example"
  description      = "example of creating tf share"
  share_proto      = "NFS"
  share_type       = "default_share_type"
  size             = 1
  share_network_id = vkcs_sharedfilesystem_sharenetwork.data.id
}
