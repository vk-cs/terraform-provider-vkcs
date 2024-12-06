resource "vkcs_kubernetes_node_group" "default_ng" {
  cluster_id = vkcs_kubernetes_cluster.k8s_cluster.id

  node_count = 1
  name       = "default"
  max_nodes  = 5
  min_nodes  = 1

  labels {
    key   = "env"
    value = "test"
  }

  labels {
    key   = "disktype"
    value = "ssd"
  }

  taints {
    key    = "taintkey1"
    value  = "taintvalue1"
    effect = "PreferNoSchedule"
  }

  taints {
    key    = "taintkey2"
    value  = "taintvalue2"
    effect = "PreferNoSchedule"
  }
}
