resource "vkcs_kubernetes_node_group_v2" "worker" {
  cluster_id = vkcs_kubernetes_cluster_v2.cluster.id
  name       = "workers"
  
  node_flavor       = data.vkcs_compute_flavor.worker.id
  availability_zone = "GZ1"
  
  scale_type = "fixed_scale"
  node_count = 1
  
  labels = {
    app         = "web-service"
  }
  
  dynamic "taints" {
    for_each = [
      {
        key    = "dedicated"
        value  = "gpu"
        effect = "NoSchedule"
      },
      {
        key    = "critical"
        value  = "true"
        effect = "NoExecute"
      }
    ]
    
    content {
      key    = taints.value.key
      value  = taints.value.value
      effect = taints.value.effect
    }
  }
  
  parallel_upgrade_chunk = 20
  
  disk_type = "high-iops"
  disk_size = 20
}
