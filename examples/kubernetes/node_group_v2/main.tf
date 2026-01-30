resource "vkcs_kubernetes_node_group_v2" "k8s_node_group" {
  cluster_id = vkcs_kubernetes_cluster_v2.k8s_cluster.id
  name       = "k8s-node-group"

  node_flavor       = data.vkcs_compute_flavor.basic.id
  availability_zone = "MS1"

  scale_type             = "fixed_scale"
  fixed_scale_node_count = 3

  parallel_upgrade_chunk = 40

  disk_type = "high-iops"
  disk_size = 20

  depends_on = [
    vkcs_kubernetes_cluster_v2.k8s_cluster,
  ]
}
