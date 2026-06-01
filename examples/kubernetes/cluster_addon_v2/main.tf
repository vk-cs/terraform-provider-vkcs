resource "vkcs_kubernetes_cluster_addon_v2" "ingress_nginx" {
  cluster_id           = vkcs_kubernetes_cluster_v2.k8s_cluster.id
  addon_id             = data.vkcs_kubernetes_addon_v2.ingress_nginx.addon_id
  addon_version_id     = data.vkcs_kubernetes_addon_v2.ingress_nginx.id
  namespace            = "ingress-nginx"
  values               = data.vkcs_kubernetes_addon_v2.ingress_nginx.values_template
  addon_name           = "ingress-nginx"

  depends_on = [
    vkcs_kubernetes_node_group_v2.k8s_node_group
  ]
}
