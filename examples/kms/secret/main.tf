resource "vkcs_kms_secret" "secret" {
  path              = "credentials"
  delete_protection = false

  data_json = jsonencode({
    username = "my-username"
    password = "my-password"
  })
}
