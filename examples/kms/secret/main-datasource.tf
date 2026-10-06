data "vkcs_kms_secret" "secret" {
  path = "application/config"
}

output "secret_data" {
  value     = jsondecode(data.vkcs_kms_secret.secret.data_json)
  sensitive = true
}
