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
