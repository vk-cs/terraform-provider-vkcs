resource "vkcs_kms_secret" "secrets_list_secret_0" {
  path              = "secrets_list_secret_0"
  delete_protection = false
  data_json = jsonencode({
    username = "someone"
  })
}

resource "vkcs_kms_secret" "secrets_list_secret_1" {
  path              = "secrets_list_secret_1"
  delete_protection = false
  data_json = jsonencode({
    username = "someone-else"
  })
}

data "vkcs_kms_secrets_list" "secrets" {
  depends_on = [ vkcs_kms_secret.secrets_list_secret_0, vkcs_kms_secret.secrets_list_secret_1 ]
}

output "kms_secrets" {
  value = data.vkcs_kms_secrets_list.secrets.secrets
}
