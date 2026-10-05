data "vkcs_kms_key_encrypt" "encrypted" {
  key       = "key-tf-example"
  plaintext = base64encode("example message")
}

output "ciphertext" {
  value = data.vkcs_kms_key_encrypt.encrypted.ciphertext
}
