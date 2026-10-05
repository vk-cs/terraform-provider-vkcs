resource "vkcs_kms_secret" "secret" {
  path = "application/config"

  data = {
    username = "example-user"
  }
}
