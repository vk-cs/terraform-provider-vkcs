data "vkcs_kubernetes_addons" "cluster-addons" {
  cluster_id = vkcs_kubernetes_cluster.k8s-cluster.id
}
