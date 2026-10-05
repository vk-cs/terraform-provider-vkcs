data "vkcs_kms_key" "key" {
  name = "key-tf-example"
}

output "key_type" {
  value = data.vkcs_kms_key.key.type
}
