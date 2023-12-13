data "vkcs_blockstorage_volume" "data" {
  name = "data-tf-example"
  # This is unnecessary in real life.
  # This is required here to let the example work with volume resource example.
  depends_on  = [vkcs_blockstorage_volume.data]
}
