resource "vkcs_kms_key" "encrypt_decrypt_key" {
  name             = "encrypt_decrypt_key"
  type             = "aes128-gcm96"
  deletion_allowed = true
}
