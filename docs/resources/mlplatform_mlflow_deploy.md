---
subcategory: "ML Platform"
layout: "vkcs"
page_title: "vkcs: vkcs_mlplatform_mlflow_deploy"
description: |-
  Manages a ML Platform MLFlow Deploy instance resource within VKCS.
---

# vkcs_mlplatform_mlflow_deploy

Manages a ML Platform Deploy resource.

**New since v0.6.0**.

## Example Usage
```terraform
resource "vkcs_mlplatform_mlflow_deploy" "deploy" {
  name               = "tf-example"
  flavor_id          = data.vkcs_compute_flavor.basic.id
  mlflow_instance_id = vkcs_mlplatform_mlflow.mlflow.id
  availability_zone  = "GZ1"
  boot_volume = {
    size        = 50
    volume_type = "ceph-ssd"
  }

  data_volumes = [
    {
      size        = 60
      volume_type = "ceph-ssd"
    },
  ]
  networks = [
    {
      network_id = vkcs_networking_network.app.id
    },
  ]
}
```

## Argument Reference
- `availability_zone` **required** *string* &rarr;  The availability zone in which to create the resource. Changing this creates a new resource

- `boot_volume` ***required*** &rarr;  Instance's boot volume configuration
  - `volume_type` **required** *string* &rarr;  Type of the volume

  - `size` optional *number* &rarr;  Size of the volume

  - `name` read-only *string* &rarr;  Name of the volume

  - `volume_id` read-only *string* &rarr;  ID of the volume


- `flavor_id` **required** *string* &rarr;  Flavor ID

- `mlflow_instance_id` **required** *string* &rarr;  MLFlow instance ID

- `name` **required** *string* &rarr;  Instance name. Changing this creates a new resource

- `networks`  *list* &rarr;  Network configuration
  - `network_id` **required** *string* &rarr;  ID of the network

  - `ip_pool` optional *string* &rarr;  ID of the ip pool


- `data_volumes`  *list* &rarr;  Instance's data volumes configuration
  - `size` **required** *number* &rarr;  Size of the volume

  - `volume_type` **required** *string* &rarr;  Type of the volume

  - `name` read-only *string* &rarr;  Name of the volume

  - `volume_id` read-only *string* &rarr;  ID of the volume


- `region` optional *string* &rarr;  The `region` in which ML Platform client is obtained, defaults to the provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Creation timestamp

- `dns_name` *string* &rarr;  DNS name

- `id` *string* &rarr;  ID of the resource

- `private_ip` *string* &rarr;  Private IP address



## Import

ML Platform JupyterHub instance can be imported using the `id`, e.g.
```shell
terraform import vkcs_mlplatform_mlflow_deploy.mymlflowdeploy 0cade671-81b5-43c5-83e1-2a659378d53a
```
