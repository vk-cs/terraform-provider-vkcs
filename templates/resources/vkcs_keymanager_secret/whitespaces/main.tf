resource "vkcs_keymanager_secret" "secret_1" {
  name                     = "password"
  payload                  = "${base64encode("password with the whitespace at the end ")}"
  secret_type              = "passphrase"
  payload_content_type     = "application/octet-stream"
  payload_content_encoding = "base64"
}
