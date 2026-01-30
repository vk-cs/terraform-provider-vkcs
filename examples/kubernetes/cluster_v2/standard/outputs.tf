output "cluster_id" {
  description = "ID of Kubernetes cluster"
  value       = vkcs_kubernetes_cluster_v2.cluster.id
}

output "cluster_uuid" {
  description = "UUID of Kubernetes cluster"
  value       = vkcs_kubernetes_cluster_v2.cluster.uuid
}

output "cluster_name" {
  description = "Name of Kubernetes cluster"
  value       = vkcs_kubernetes_cluster_v2.cluster.name
}

output "cluster_type" {
  description = "Type of cluster"
  value       = vkcs_kubernetes_cluster_v2.cluster.cluster_type
}

output "region" {
  description = "Region where cluster is deployed"
  value       = vkcs_kubernetes_cluster_v2.cluster.region
}

output "cluster_version" {
  description = "Kubernetes version"
  value       = vkcs_kubernetes_cluster_v2.cluster.version
}

output "cluster_description" {
  description = "Description of Kubernetes cluster"
  value       = vkcs_kubernetes_cluster_v2.cluster.description
}

output "availability_zones" {
  description = "Availability zones of cluster"
  value       = vkcs_kubernetes_cluster_v2.cluster.availability_zones
}

output "labels" {
  description = "Labels of cluster"
  value       = vkcs_kubernetes_cluster_v2.cluster.labels
}

output "master_count" {
  description = "Number of master nodes"
  value       = vkcs_kubernetes_cluster_v2.cluster.master_count
}

output "master_flavor" {
  description = "ID of master flavor"
  value       = vkcs_kubernetes_cluster_v2.cluster.master_flavor
}

output "network_id" {
  description = "ID of cluster network"
  value       = vkcs_kubernetes_cluster_v2.cluster.network_id
}

output "subnet_id" {
  description = "ID of cluster subnet"
  value       = vkcs_kubernetes_cluster_v2.cluster.subnet_id
}

output "cluster_has_public_ip" {
  description = "Does cluster have public IP"
  value       = vkcs_kubernetes_cluster_v2.cluster.enable_public_ip
}

output "external_network_id" {
  description = "ID of external network"
  value       = vkcs_kubernetes_cluster_v2.cluster.external_network_id
}

output "cluster_insecure_registries" {
  description = "List of cluster insecure registries"
  value       = vkcs_kubernetes_cluster_v2.cluster.insecure_registries
}

output "network_plugin" {
  description = "Network plugin used in cluster"
  value       = vkcs_kubernetes_cluster_v2.cluster.network_plugin
}

output "pods_ipv4_cidr" {
  description = "IPv4 CIDR for pods"
  value       = vkcs_kubernetes_cluster_v2.cluster.pods_ipv4_cidr
}

output "loadbalancer_subnet_id" {
  description = "ID of loadbalancer subnet"
  value       = vkcs_kubernetes_cluster_v2.cluster.loadbalancer_subnet_id
}

output "loadbalancer_allowed_cidrs" {
  description = "List of CIDR blocks allowed to access load balancer"
  value       = vkcs_kubernetes_cluster_v2.cluster.loadbalancer_allowed_cidrs
}

output "cluster_kubeconfig" {
  description = "Cluster kubeconfig"
  value       = vkcs_kubernetes_cluster_v2.cluster.k8s_config
}

output "master_disks" {
  description = "List of master disks"
  value       = vkcs_kubernetes_cluster_v2.cluster.master_disks
}

output "created_at" {
  description = "Time at which cluster was created"
  value       = vkcs_kubernetes_cluster_v2.cluster.created_at
}

output "cluster_status" {
  description = "Cluster current status"
  value       = vkcs_kubernetes_cluster_v2.cluster.status
}

output "user_project_id" {
  description = "User project ID"
  value       = vkcs_kubernetes_cluster_v2.cluster.project_id
}

output "api_lb_fip" {
  description = "API LoadBalancer FIP (Floating IP). IP address field"
  value       = vkcs_kubernetes_cluster_v2.cluster.api_lb_fip
}

output "api_lb_vip" {
  description = "API LoadBalancer VIP (Virtual IP). IP address field"
  value       = vkcs_kubernetes_cluster_v2.cluster.api_lb_vip
}

output "api_address" {
  description = "URL address of cluster kubeapi-server"
  value       = vkcs_kubernetes_cluster_v2.cluster.api_address
}

output "cluster_node_groups" {
  description = "List of cluster node groups"
  value       = vkcs_kubernetes_cluster_v2.cluster.node_groups
}
