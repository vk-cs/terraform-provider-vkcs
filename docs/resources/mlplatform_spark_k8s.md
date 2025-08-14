---
subcategory: "ML Platform"
layout: "vkcs"
page_title: "vkcs: vkcs_mlplatform_spark_k8s"
description: |-
  Manages a ML Platform Spark K8S cluster resource within VKCS.
---

# vkcs_mlplatform_spark_k8s

Manages a ML Platform Spark K8S cluster resource.

**New since v0.7.0**.

## Example Usage
```terraform
locals {
  spark_configuration = {
    "spark.eventLog.dir"     = "s3a://spark-bucket"
    "spark.eventLog.enabled" = "true"
  }
  spark_environment_variables = {
    "ENV_VAR_1" : "env_var_1_value"
    "ENV_VAR_2" : "env_var_2_value"
  }
}

resource "vkcs_mlplatform_spark_k8s" "spark_k8s" {
  name              = "tf-example"
  availability_zone = "GZ1"
  network_id        = vkcs_networking_network.app.id
  subnet_id         = vkcs_networking_subnet.app.id

  node_groups = [
    {
      node_count          = 2
      flavor_id           = data.vkcs_compute_flavor.basic.id
      autoscaling_enabled = true
      min_nodes           = 2
      max_nodes           = 100
    }
  ]
  cluster_mode = "DEV"
  registry_id  = vkcs_mlplatform_k8s_registry.k8s_registry.id
  ip_pool      = data.vkcs_networking_network.extnet.id

  suspend_after_inactive_min = 120
  delete_after_inactive_min  = 1440

  spark_configuration   = yamlencode(local.spark_configuration)
  environment_variables = yamlencode(local.spark_environment_variables)
}
```

## Argument Reference
- `availability_zone` **required** *string* &rarr;  The availability zone in which to create the resource. Changing this creates a new resource

- `cluster_mode` **required** *string* &rarr;  Cluster Mode. Should be `DEV` or `PROD`. Changing this creates a new resource

- `ip_pool` **required** *string* &rarr;  ID of the ip pool. Changing this creates a new resource

- `name` **required** *string* &rarr;  Cluster name. Changing this creates a new resource

- `network_id` **required** *string* &rarr;  ID of the network. Changing this creates a new resource

- `node_groups`  *list* &rarr;  Cluster's node groups configuration
    - `autoscaling_enabled` **required** *boolean* &rarr;  Enables autoscaling for node group

    - `flavor_id` **required** *string* &rarr;  ID of the flavor to be used in nodes

    - `max_nodes` optional *number* &rarr;  Maximum number of nodes in node group. It is used only when autoscaling is enabled

    - `min_nodes` optional *number* &rarr;  Minimum count of nodes in node group. It is used only when autoscaling is enabled

    - `node_count` optional *number* &rarr;  Count of nodes in node group


- `registry_id` **required** *string* &rarr;  ID of the K8S registry to use with cluster. Changing this creates a new resource

- `delete_after_inactive_min` optional *number* &rarr;  Timeout of cluster inactivity before deletion, in minutes. Changing this creates a new resource

- `environment_variables` optional *string* &rarr;  Environment variables. Read more about this parameter here: https://cloud.vk.com/docs/en/ml/spark-to-k8s/instructions/create. Changing this creates a new resource

- `region` optional *string* &rarr;  The `region` in which ML Platform client is obtained, defaults to the provider's `region`.

- `spark_configuration` optional *string* &rarr;  Spark configuration. Read more about this parameter here: https://cloud.vk.com/docs/en/ml/spark-to-k8s/instructions/create. Changing this creates a new resource

- `subnet_id` optional *string* &rarr;  ID of the subnet. Changing this creates a new resource

- `suspend_after_inactive_min` optional *number* &rarr;  Timeout of cluster inactivity before suspending, in minutes. Changing this creates a new resource


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `control_instance_id` *string* &rarr;  ID of the control instance

- `history_server_url` *string* &rarr;  URL of the history server

- `id` *string* &rarr;  ID of the resource

- `inactive_min` *number* &rarr;  Current time of cluster inactivity, in minutes

- `s3_bucket_name` *string* &rarr;  S3 bucket name



## Import

ML Platform Spark K8S cluster can be imported using the `id`, e.g.
```shell
terraform import vkcs_mlplatform_spark_k8s.mysparkk8s 32cc47a5-9726-454f-bffa-6723f21fbbc7
```
