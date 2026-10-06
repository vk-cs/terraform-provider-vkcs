resource "vkcs_kms_secret" "secret" {
  path = "application/config"

  data_json = jsonencode({
    username = "example-user"
  })
}
