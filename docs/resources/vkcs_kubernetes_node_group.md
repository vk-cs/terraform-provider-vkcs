---
layout: "vkcs"
page_title: "vkcs: vkcs_kubernetes_node_group"
description: |-
  Manages clusters node group.
---

# vkcs_kubernetes_node_group

Provides a cluster node group resource. This can be used to create, modify and delete cluster's node group.

## Example Usage
```terraform
data "vkcs_kubernetes_clustertemplate" "ct" {
    version = "1.24"
}

resource "vkcs_kubernetes_cluster" "k8s-cluster" {
    depends_on = [
        vkcs_networking_router_interface.k8s,
    ]

    name                = "k8s-cluster"
    cluster_template_id = data.vkcs_kubernetes_clustertemplate.ct.id
    master_flavor       = data.vkcs_compute_flavor.k8s.id
    master_count        = 1

    network_id          = vkcs_networking_network.k8s.id
    subnet_id           = vkcs_networking_subnet.k8s-subnetwork.id
    floating_ip_enabled = true
    availability_zone   = "MS1"
    insecure_registries = ["1.2.3.4"]
}

resource "vkcs_kubernetes_node_group" "default_ng" {
    cluster_id = vkcs_kubernetes_cluster.k8s-cluster.id

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
```
## Argument Reference
- `cluster_id` **String** (***Required***) The UUID of the existing cluster.

- `name` **String** (***Required***) The name of node group to create. Changing this will force to create a new node group.

- `node_count` **Number** (***Required***) The node count for this node group. Should be greater than 0. If `autoscaling_enabled` parameter is set, this attribute will be ignored during update.

- `autoscaling_enabled` **Boolean** (*Optional*) Determines whether the autoscaling is enabled.

- `availability_zones` **String** (*Optional*) The list of availability zones of the node group. Zones `MS1` and  `GZ1` are available. By default, node group is being created at cluster's zone.
**Important:** Receiving default AZ add it manually to your main.tf config to sync it with state to avoid node groups force recreation in the future.

- `flavor_id` **String** (*Optional*) The flavor UUID of this node group.

- `labels` (*Optional*) The list of objects representing representing additional properties of the node group. Each object should have attribute "key". Object may also have optional attribute "value".
  - `key` **String** (***Required***)

  - `value` **String** (*Optional*)

- `max_node_unavailable` **Number** (*Optional*) The maximum number of nodes that can fail during an upgrade. The default value is 25 percent.

- `max_nodes` **Number** (*Optional*) The maximum allowed nodes for this node group.

- `min_nodes` **Number** (*Optional*) The minimum allowed nodes for this node group. Default to 0 if not set.

- `taints` (*Optional*) The list of objects representing node group taints. Each object should have following attributes: key, value, effect.
  - `effect` **String** (***Required***)

  - `key` **String** (***Required***)

  - `value` **String** (***Required***)

- `volume_size` **Number** (*Optional*) The size in GB for volume to load nodes from. Changing this will force to create a new node group.

- `volume_type` **String** (*Optional*) The volume type to load nodes from. Changing this will force to create a new node group.


## Attributes Reference
- `cluster_id` **String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `node_count` **Number** See Argument Reference above.

- `autoscaling_enabled` **Boolean** See Argument Reference above.

- `availability_zones` **String** See Argument Reference above.

- `flavor_id` **String** See Argument Reference above.

- `labels`  See Argument Reference above.
  - `key` **String**

  - `value` **String**

- `max_node_unavailable` **Number** See Argument Reference above.

- `max_nodes` **Number** See Argument Reference above.

- `min_nodes` **Number** See Argument Reference above.

- `taints`  See Argument Reference above.
  - `effect` **String**

  - `key` **String**

  - `value` **String**

- `volume_size` **Number** See Argument Reference above.

- `volume_type` **String** See Argument Reference above.

- `created_at` **String** The time at which node group was created.

- `id` **String** ID of the resource.

- `state` **String** Determines current state of node group (RUNNING, SHUTOFF, ERROR).

- `updated_at` **String** The time at which node group was created.

- `uuid` **String** The UUID of the cluster's node group.



## Import

Node groups can be imported using the `id`, e.g.

```shell
terraform import vkcs_kubernetes_node_group.ng aa14de9c-c5f5-4cc0-869c-ce655419df76
```
