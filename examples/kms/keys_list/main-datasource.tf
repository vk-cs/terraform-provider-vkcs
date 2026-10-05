data "vkcs_kms_keys_list" "keys" {}

output "kms_keys" {
  value = data.vkcs_kms_keys_list.keys.keys
}
