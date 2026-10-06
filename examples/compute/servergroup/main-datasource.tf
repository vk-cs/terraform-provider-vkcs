data "vkcs_compute_servergroup" "example" {
  name = "custom-group-tf-example"
  # This is unnecessary in real life.
  # This is required here to let the example work with the server group resource example.
  depends_on = [vkcs_compute_servergroup.cusom_group]
}
