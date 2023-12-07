data "vkcs_compute_availability_zones" "zones" {}

output "available_zones" {
  value = data.vkcs_compute_availability_zones.zones.names
}
