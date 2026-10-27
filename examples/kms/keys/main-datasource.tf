data "vkcs_kms_keys" "keys" {
  depends_on = [vkcs_kms_key.keys_key_0, vkcs_kms_key.keys_key_1]
}

output "kms_keys" {
  value = data.vkcs_kms_keys.keys.keys
}
