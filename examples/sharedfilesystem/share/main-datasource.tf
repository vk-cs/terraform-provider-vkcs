data "vkcs_sharedfilesystem_share" "data" {
  name = "share-data-tf-example"
  share_network_id = vkcs_sharedfilesystem_sharenetwork.data.id
  # This is unnecessary in real life.
  # This is required here to let the example work with share resource example. 
  depends_on = [ vkcs_sharedfilesystem_share.data ]
}
