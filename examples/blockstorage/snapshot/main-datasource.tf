data "vkcs_blockstorage_snapshot" "recent_snapshot" {
  name        = "snapshot-tf-example"
  most_recent = true
  # This is unnecessary in real life.
  # This is required here to let the example work with snapshot resource example.
  depends_on  = [vkcs_blockstorage_snapshot.recent_snapshot]
}
