data "vkcs_kms_key_encrypt" "encrypted" {
  key       = "key-tf-example"
  plaintext = base64encode("example message")
}

data "vkcs_kms_key_decrypt" "decrypted" {
  key        = "key-tf-example"
  ciphertext = data.vkcs_kms_key_encrypt.encrypted.ciphertext
}

output "message" {
  value     = base64decode(data.vkcs_kms_key_decrypt.decrypted.plaintext)
  sensitive = true
}
