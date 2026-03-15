output "node_group_id" {
  description = "ID of the node group"
  value       = vkcs_kubernetes_node_group_v2.node_group.id
}

output "node_group_uuid" {
  description = "UUID of the node group"
  value       = vkcs_kubernetes_node_group_v2.node_group.uuid
}

output "node_group_cluster_id" {
  description = "ID of the target cluster"
  value       = vkcs_kubernetes_node_group_v2.node_group.cluster_id
}

output "node_group_name" {
  description = "Name of the node group"
  value       = vkcs_kubernetes_node_group_v2.node_group.name
}

output "node_group_node_flavor" {
  description = "Flavor ID of the nodes from node group"
  value       = vkcs_kubernetes_node_group_v2.node_group.node_flavor
}

output "node_group_availability_zone" {
  description = "Availability zone of the node group"
  value       = vkcs_kubernetes_node_group_v2.node_group.availability_zone
}

output "node_group_scale_type" {
  description = "Type of scaling for the node group. Must be either 'fixed_scale' or 'auto_scale'"
  value       = vkcs_kubernetes_node_group_v2.node_group.scale_type
}

output "node_group_node_count" {
  description = "Node count of the node group"
  value       = vkcs_kubernetes_node_group_v2.node_group.node_count
}

output "node_group_auto_scale_min_size" {
  description = "Minimum allowed nodes for this node group (only for auto_scale type)"
  value       = vkcs_kubernetes_node_group_v2.node_group.auto_scale_min_size
}

output "node_group_auto_scale_max_size" {
  description = "Maximum allowed nodes for this node group (only for auto_scale type)"
  value       = vkcs_kubernetes_node_group_v2.node_group.auto_scale_max_size
}

output "node_group_labels" {
  description = "Key-value pairs representing additional properties of the node group"
  value       = vkcs_kubernetes_node_group_v2.node_group.labels
}

output "node_group_taints" {
  description = "List of objects representing node group taints. Each object has key, value, and effect attributes"
  value       = vkcs_kubernetes_node_group_v2.node_group.taints
}

output "node_group_disk_type" {
  description = "Volume type to load nodes from"
  value       = vkcs_kubernetes_node_group_v2.node_group.disk_type
}

output "node_group_disk_size" {
  description = "Size in GB for volume to load nodes from"
  value       = vkcs_kubernetes_node_group_v2.node_group.disk_size
}

output "node_group_parallel_upgrade_chunk" {
  description = "Maximum percent of nodes that can be unavailable during an upgrade"
  value       = vkcs_kubernetes_node_group_v2.node_group.parallel_upgrade_chunk
}

output "node_group_region" {
  description = "Region used for the node group"
  value       = vkcs_kubernetes_node_group_v2.node_group.region
}

output "node_group_created_at" {
  description = "Time at which node group was created"
  value       = vkcs_kubernetes_node_group_v2.node_group.created_at
}

output "node_group_is_autoscale" {
  description = "Whether the node group uses auto-scaling"
  value       = vkcs_kubernetes_node_group_v2.node_group.scale_type == "auto_scale"
}

output "node_group_is_fixed_scale" {
  description = "Whether the node group uses fixed scaling"
  value       = vkcs_kubernetes_node_group_v2.node_group.scale_type == "fixed_scale"
}
