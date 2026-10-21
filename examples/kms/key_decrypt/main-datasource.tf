data "vkcs_kms_key_encrypt" "encrypted" {
  key       = vkcs_kms_key.encrypt_decrypt_key.id
  plaintext = base64encode("a secret message")
}

output "encrypted" {
  value = data.vkcs_kms_key_encrypt.encrypted.ciphertext
}

data "vkcs_kms_key_decrypt" "decrypted" {
  key        = vkcs_kms_key.encrypt_decrypt_key.id
  ciphertext = data.vkcs_kms_key_encrypt.encrypted.ciphertext
}

output "decrypted" {
  value     = base64decode(data.vkcs_kms_key_decrypt.decrypted.plaintext)
  sensitive = true
}
