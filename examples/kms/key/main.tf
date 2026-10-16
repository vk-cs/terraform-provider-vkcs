resource "vkcs_kms_key" "key" {
  name             = "some-key"
  type             = "aes256-gcm96"
  deletion_allowed = true
}

data "vkcs_kms_key" "key" {
  name = vkcs_kms_key.key.id
}

output "key_type" {
  value = data.vkcs_kms_key.key.type
}

output "key_deletion_allowed" {
  value = data.vkcs_kms_key.key.deletion_allowed
}
