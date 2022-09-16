resource "vkcs_keymanager_secret" "secret_1" {
  algorithm            = "aes"
  bit_length           = 256
  mode                 = "cbc"
  name                 = "mysecret"
  payload              = "foobar"
  payload_content_type = "text/plain"
  secret_type          = "passphrase"

  metadata = {
    key = "foo"
  }
}
