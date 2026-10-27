resource "vkcs_kms_key" "encrypt_key" {
  name             = "encrypt_key"
  type             = "aes128-gcm96"
  deletion_allowed = true
}
