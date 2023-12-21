---
subcategory: "ML Platform"
layout: "vkcs"
page_title: "vkcs: vkcs_mlplatform_jupyterhub"
description: |-
  Manages a ML Platform JupyterHub instance resource within VKCS.
---

# vkcs_mlplatform_jupyterhub

Manages a ML Platform JupyterHub resource.

**New since v0.6.0**.

## Example Usage
```terraform
resource "vkcs_mlplatform_jupyterhub" "jupyterhub" {
  name              = "tf-example"
  admin_name        = "admin"
  admin_password    = "Password!"
  flavor_id         = data.vkcs_compute_flavor.basic.id
  availability_zone = "GZ1"
  boot_volume = {
    volume_type = "ceph-ssd"
  }
  data_volumes = [
    {
      size        = 60
      volume_type = "ceph-ssd"
    },
    {
      size        = 70
      volume_type = "ceph-ssd"
    }
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


- `data_volumes`  *list* &rarr;  Instance's data volumes configuration
  - `size` **required** *number* &rarr;  Size of the volume

  - `volume_type` **required** *string* &rarr;  Type of the volume

  - `name` read-only *string* &rarr;  Name of the volume

  - `volume_id` read-only *string* &rarr;  ID of the volume


- `flavor_id` **required** *string* &rarr;  Flavor ID

- `name` **required** *string* &rarr;  Instance name. Changing this creates a new resource

- `networks`  *list* &rarr;  Network configuration
  - `network_id` **required** *string* &rarr;  ID of the network

  - `ip_pool` optional *string* &rarr;  ID of the ip pool


- `admin_name` optional *string* &rarr;  JupyterHub admin name. Changing this creates a new resource

- `admin_password` optional sensitive *string* &rarr;  JupyterHub admin password. Changing this creates a new resource

- `domain_name` optional *string* &rarr;  Domain name. Changing this creates a new resource

- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`.

- `s3fs_bucket` optional *string* &rarr;  Connect specified s3 bucket to instance as volume. Changing this creates a new resource


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Creation timestamp

- `dns_name` *string* &rarr;  DNS name

- `id` *string* &rarr;  ID of the resource

- `private_ip` *string* &rarr;  Private IP address



## Import

ML Platform JupyterHub instance can be imported using the `id`, e.g.
```shell
terraform import vkcs_mlplatform_jupyterhub.myjupyterhub 3a679dd9-0942-49b0-b233-95de5a5a9502
```
