data "vkcs_kubernetes_volume_types_v2" "available_volume_types" {}

output "available_volume_types" {
  description = "A set of storage volume types that can be selected as the root disk for node groups."
  value       = data.vkcs_kubernetes_volume_types_v2.available_volume_types.volume_types
}
