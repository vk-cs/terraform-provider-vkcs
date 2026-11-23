resource "vkcs_kubernetes_node_group_v2" "k8s_node_group" {
  cluster_id = vkcs_kubernetes_cluster_v2.k8s_cluster.id
  name       = "k8s-node-group"

  node_flavor       = data.vkcs_compute_flavor.worker.id
  availability_zone = "MS1"

  scale_type             = "fixed_scale"
  fixed_scale_node_count = 3

  parallel_upgrade_chunk = 40

  disk_type = "high-iops"
  disk_size = 20

  # Optional: additional security groups attached to the nodes,
  # in addition to the default security groups.
  # extra_security_groups = ["c906d669-b9ab-4e37-b34b-a2a96b5ed626"]

  depends_on = [
    vkcs_kubernetes_cluster_v2.k8s_cluster,
  ]
}
