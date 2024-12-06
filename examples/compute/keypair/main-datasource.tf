data "vkcs_compute_keypair" "generated_key" {
  name = "generated-key-tf-example"
  # This is unnecessary in real life.
  # This is required here to let the example work with keypair resource example. 
  depends_on = [vkcs_compute_keypair.generated_key]
}
