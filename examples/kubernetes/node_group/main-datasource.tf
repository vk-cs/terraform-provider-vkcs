data "vkcs_kubernetes_node_group" "k8s-node-group" {
  id = vkcs_kubernetes_node_group.default_ng.id
}
