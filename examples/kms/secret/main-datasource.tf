data "vkcs_kms_secret" "secret" {
  path = "application/config"
}

output "secret_data" {
  value     = data.vkcs_kms_secret.secret.data
  sensitive = true
}
