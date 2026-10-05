resource "vkcs_kms_key" "key" {
  name             = "key-tf-example"
  type             = "aes256-gcm96"
  deletion_allowed = false
}
