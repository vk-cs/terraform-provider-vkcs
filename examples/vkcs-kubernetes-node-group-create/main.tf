resource "vkcs_kubernetes_node_group" "default_ng" {
    cluster_id = "your_cluster_id"

    node_count = 1
    name = var.name
    max_nodes = 5
    min_nodes = 1
    max_node_unavailable = var.max-node-unavailable

    labels {
        key = "env"
        value = "test"
    }

    labels {
        key = "disktype"
        value = "ssd"
    }
    
    taints {
        key = "taintkey1"
        value = "taintvalue1"
        effect = "PreferNoSchedule"
    }

    taints {
        key = "taintkey2"
        value = "taintvalue2"
        effect = "PreferNoSchedule"
    }
}
