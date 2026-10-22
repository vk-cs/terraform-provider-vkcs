resource "vkcs_kms_secret" "secrets_secret_0" {
  path              = "secrets_secret_0"
  delete_protection = false
  data_json = jsonencode({
    username = "someone"
  })
}

resource "vkcs_kms_secret" "secrets_secret_1" {
  path              = "secrets_secret_1"
  delete_protection = false
  data_json = jsonencode({
    username = "someone-else"
  })
}
