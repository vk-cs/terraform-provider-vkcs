terraform {
  required_providers {
    mcs = {
      source = "vk-cs/vkcs"
      version = "~> 0.1.0"
    }
  }
}

data "vkcs_kubernetes_cluster" "your_cluster" {
  cluster_id = "your_cluster_uuid"
}

resource "vkcs_kubernetes_node_group" "default_ng" {
  cluster_id = data.vkcs_kubernetes_cluster.your_cluster.id

  node_count = 1
  name = "default"
  max_nodes = 5
  min_nodes = 1

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
