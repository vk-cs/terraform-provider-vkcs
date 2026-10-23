resource "vkcs_kms_key" "key" {
  name             = "some-key"
  type             = "aes256-gcm96"
  deletion_allowed = true
}
