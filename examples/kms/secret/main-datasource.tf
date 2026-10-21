data "vkcs_kms_secret" "secret" {
  path = vkcs_kms_secret.secret.id
}

output "password" {
  value     = jsondecode(data.vkcs_kms_secret.secret.data_json).password
  sensitive = true
}
