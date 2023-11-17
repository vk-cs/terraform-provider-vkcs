resource "vkcs_keymanager_secret" "priv_key" {
  name                 = "priv-key-tf-example"
  secret_type          = "private"
  payload_content_type = "text/plain"
  payload              = file("${path.module}/private-key.key")
}
