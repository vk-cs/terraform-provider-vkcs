data "vkcs_kms_key" "key" {
  name = vkcs_kms_key.key.id
}

output "key_deletion_allowed" {
  value = data.vkcs_kms_key.key.deletion_allowed
}
