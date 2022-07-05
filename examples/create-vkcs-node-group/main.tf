terraform {
    required_providers {
        vkcs = {
            source  = "vk-cs/vkcs"
            version = "~> 0.1.0"
        }
    }
}

data "vkcs_kubernetes_cluster" "your_cluster" {
    cluster_id = "1b322180-503b-44e3-8d92-934b6e574e66"
}

resource "vkcs_kubernetes_node_group" "default_ng" {
    cluster_id = data.vkcs_kubernetes_cluster.your_cluster.id

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
