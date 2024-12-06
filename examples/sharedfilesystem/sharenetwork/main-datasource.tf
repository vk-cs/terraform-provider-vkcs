data "vkcs_sharedfilesystem_sharenetwork" "data" {
  name = "sharenetwork-tf-example"
  # This is unnecessary in real life.
  # This is required here to let the example work with sharenetwork resource example. 
  depends_on = [vkcs_sharedfilesystem_sharenetwork.data]
}
