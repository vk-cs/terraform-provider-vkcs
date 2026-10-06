resource "vkcs_kms_secret" "secret" {
  path              = "application/config"
  delete_protection = true

  data_json = jsonencode({
    username = "example-user"
  })
}
