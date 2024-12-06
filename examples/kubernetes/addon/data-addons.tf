data "vkcs_kubernetes_addons" "cluster_addons" {
  cluster_id = vkcs_kubernetes_cluster.k8s_cluster.id
}
