data "vkcs_kms_secrets_list" "secrets" {}

output "kms_secrets" {
  value = data.vkcs_kms_secrets_list.secrets.secrets
}
