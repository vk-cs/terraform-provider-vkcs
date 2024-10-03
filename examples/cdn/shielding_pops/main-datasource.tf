data "vkcs_cdn_shielding_pops" "pops" {}

output "shielding_locations" {
  value = data.vkcs_cdn_shielding_pops.pops.shielding_pops
}
