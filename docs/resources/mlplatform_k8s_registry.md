---
subcategory: "ML Platform"
layout: "vkcs"
page_title: "vkcs: vkcs_mlplatform_k8s_registry"
description: |-
  Manages a ML Platform K8S Registry instance resource within VKCS.
---

# vkcs_mlplatform_k8s_registry

Manages a ML Platform K8SRegistry resource.

## Example Usage
```terraform
resource "vkcs_mlplatform_k8s_registry" "k8s_registry" {
  name              = "tf-example"
  admin_name        = "admin"
  admin_password    = "dM8Ao21,0S264iZp"
  flavor_id         = data.vkcs_compute_flavor.basic.id
  availability_zone = "GZ1"
  boot_volume = {
    volume_type = "ceph-ssd"
  }
  networks = [
    {
      network_id = vkcs_networking_network.app.id
      ip_pool = data.vkcs_networking_network.extnet.id
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

- `name` **required** *string* &rarr;  Instance name. Changing this creates a new resource

- `networks`  *list* &rarr;  Network configuration
  - `ip_pool` **required** *string* &rarr;  ID of the ip pool

  - `network_id` **required** *string* &rarr;  ID of the network


- `admin_name` optional *string* &rarr;  K8SRegistry admin name. Changing this creates a new resource

- `admin_password` optional sensitive *string* &rarr;  K8SRegistry admin password. Changing this creates a new resource

- `data_volumes`  *list* &rarr;  Instance's data volumes configuration
  - `size` **required** *number* &rarr;  Size of the volume

  - `volume_type` **required** *string* &rarr;  Type of the volume

  - `name` read-only *string* &rarr;  Name of the volume

  - `volume_id` read-only *string* &rarr;  ID of the volume


- `domain_name` optional *string* &rarr;  Domain name. Changing this creates a new resource

- `region` optional *string* &rarr;  The `region` in which ML Platform client is obtained, defaults to the provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `created_at` *string* &rarr;  Creation timestamp

- `dns_name` *string* &rarr;  DNS name

- `id` *string* &rarr;  ID of the resource

- `private_ip` *string* &rarr;  Private IP address



## Import

ML Platform K8S Registry instance can be imported using the `id`, e.g.
```shell
terraform import vkcs_mlplatform_k8s_registry.myk8sregistry 0229eb40-5b56-4ab1-857f-453848a542f3
```
