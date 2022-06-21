resource "vkcs_compute_keypair" "keypair" {
  name = "test-keypair"
}

output "public_key" {
  value = vkcs_compute_keypair.keypair.public_key
}

output "private_key" {
  value = vkcs_compute_keypair.keypair.private_key
  sensitive = true
}
