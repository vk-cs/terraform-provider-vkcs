resource "vkcs_compute_keypair" "generated-key" {
  name = "generated-key-tf-example"
}

output "public_key" {
  value = vkcs_compute_keypair.generated-key.public_key
}

output "private_key" {
  value = vkcs_compute_keypair.generated-key.private_key
  sensitive = true
}
