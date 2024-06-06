data "vkcs_kubernetes_addon" "kube-prometheus-stack" {
  cluster_id = vkcs_kubernetes_cluster.k8s-cluster.id
  name       = "kube-prometheus-stack"
  version    = "54.2.2"
}
