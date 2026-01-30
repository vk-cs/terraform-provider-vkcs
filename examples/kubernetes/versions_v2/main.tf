data "vkcs_kubernetes_versions_v2" "available_versions" {}

output "available_kubernetes_versions" {
  description = "A set of Kubernetes versions that can be used to deploy a new cluster."
  value       = data.vkcs_kubernetes_versions_v2.available_versions.k8s_versions
}
