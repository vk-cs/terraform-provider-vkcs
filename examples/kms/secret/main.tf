resource "vkcs_kms_secret" "secret" {
  path              = "credentials"
  delete_protection = false

  data_json = jsonencode({
    username = "my-username"
    password = "my-password"
  })
}

data "vkcs_kms_secret" "secret" {
  path = vkcs_kms_secret.secret.id
}

output "secret_data" {
  value     = jsondecode(data.vkcs_kms_secret.secret.data_json)
  sensitive = true
}
