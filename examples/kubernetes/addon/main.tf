resource "vkcs_kubernetes_addon" "kube-prometheus-stack" {
  cluster_id           = vkcs_kubernetes_cluster.k8s-cluster.id
  addon_id             = data.vkcs_kubernetes_addon.kube-prometheus-stack.id
  namespace            = "monitoring"
  configuration_values = data.vkcs_kubernetes_addon.kube-prometheus-stack.configuration_values

  depends_on = [
    vkcs_kubernetes_node_group.default_ng
  ]
}
