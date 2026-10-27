data "vkcs_kms_key_encrypt" "encrypted" {
  key       = vkcs_kms_key.encrypt_key.id
  plaintext = base64encode("a secret message")
}

output "encrypted" {
  value = data.vkcs_kms_key_encrypt.encrypted.ciphertext
}
