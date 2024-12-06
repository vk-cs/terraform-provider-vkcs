data "vkcs_kubernetes_addon" "kube_prometheus_stack" {
  cluster_id = vkcs_kubernetes_cluster.k8s_cluster.id
  name       = "kube-prometheus-stack"
  version    = "54.2.2"
}
