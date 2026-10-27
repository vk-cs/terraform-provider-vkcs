data "vkcs_kms_secrets" "secrets" {
  depends_on = [vkcs_kms_secret.secrets_secret_0, vkcs_kms_secret.secrets_secret_1]
}

output "kms_secrets" {
  value = data.vkcs_kms_secrets.secrets.secrets
}
