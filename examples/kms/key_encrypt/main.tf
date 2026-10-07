resource "vkcs_kms_key" "encrypt_key" {
  name             = "encrypt_key"
  type             = "aes128-gcm96"
  deletion_allowed = true
}

data "vkcs_kms_key_encrypt" "encrypted" {
  key       = vkcs_kms_key.encrypt_key.id
  plaintext = base64encode("a secret message")
}

output "encrypted" {
  value = data.vkcs_kms_key_encrypt.encrypted.ciphertext
}
