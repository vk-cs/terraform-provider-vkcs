resource "vkcs_kms_key" "keys_key_0" {
  name             = "keys_key_0"
  type             = "aes128-gcm96"
  deletion_allowed = true
}

resource "vkcs_kms_key" "keys_key_1" {
  name             = "keys_key_1"
  deletion_allowed = true
}
