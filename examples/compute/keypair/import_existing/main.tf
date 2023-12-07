resource "vkcs_compute_keypair" "existing-key" {
  name       = "existing-key-tf-example"
  public_key = file("${path.module}/public_key.key")
}
