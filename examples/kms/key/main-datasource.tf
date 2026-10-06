data "vkcs_kms_key" "key" {
  name = "key-tf-example"
}

output "key_type" {
  value = data.vkcs_kms_key.key.type
}

output "key_deletion_allowed" {
  value = data.vkcs_kms_key.key.deletion_allowed
}
