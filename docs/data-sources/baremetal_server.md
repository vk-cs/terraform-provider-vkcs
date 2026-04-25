---
subcategory: "Baremetal"
layout: "vkcs"
page_title: "vkcs: vkcs_baremetal_server"
description: |-
  Get information on a VKCS bare metal server.
---

# vkcs_baremetal_server

Use this data source to get information about a VKCS bare metal server.

## Example Usage

```terraform
data "vkcs_baremetal_server" "server" {
  id = vkcs_baremetal_server.server.id
}

output "server_output" {
  value = {
    id                = data.vkcs_baremetal_server.server.id
    name              = data.vkcs_baremetal_server.server.name
    region            = data.vkcs_baremetal_server.server.region
    availability_zone = data.vkcs_baremetal_server.server.availability_zone

    cpu_cores = data.vkcs_baremetal_server.server.cpu_cores
    cpu_types = data.vkcs_baremetal_server.server.cpu_types

    is_locked  = data.vkcs_baremetal_server.server.is_locked

    local_disk_sizes = data.vkcs_baremetal_server.server.local_disk_sizes
    ram_megabytes    = data.vkcs_baremetal_server.server.ram_megabytes

    power_state = data.vkcs_baremetal_server.server.power_state
    status      = data.vkcs_baremetal_server.server.status

    tags = data.vkcs_baremetal_server.server.tags

    image_id   = data.vkcs_baremetal_server.server.image_id
    image_name = data.vkcs_baremetal_server.server.image_name
    os_type    = data.vkcs_baremetal_server.server.os_type

    raid_type = data.vkcs_baremetal_server.server.raid_type
    flavor_id = data.vkcs_baremetal_server.server.flavor_id

    target_boot_order = data.vkcs_baremetal_server.server.target_boot_order

    local_disks_info = data.vkcs_baremetal_server.server.local_disks_info
  }
}

output "boot_devices" {
  value = [
    for b in data.vkcs_baremetal_server.server.target_boot_order :
    b.boot_device_type
  ]
}

output "disks" {
  value = [
    for d in data.vkcs_baremetal_server.server.local_disks_info : {
      path  = d.path
      size  = d.size
      type  = d.type
      model = d.model
    }
  ]
}
```

## Argument Reference
- `id` **required** *string* &rarr;  The UUID of the bare metal server.

- `region` optional *string* &rarr;  The region to fetch the bare metal server from, defaults to the provider's region.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `availability_zone` *string* &rarr;  The availability zone of this server.

- `cpu_cores` *number* &rarr;  CPU core count including hyper-threading.

- `cpu_types` *string* &rarr;  Server CPU type.

- `flavor_id` *string* &rarr;  Flavor ID of the server.

- `image_id` *string* &rarr;  The image ID used to create the server.

- `image_name` *string* &rarr;  The image name used to create the server.

- `is_locked` *boolean* &rarr;  Shows whether the server is protected. The server cannot be deleted while this flag is set.

- `local_disk_sizes` *number* &rarr;  Local disk sizes in gigabytes.

- `local_disks_info`  *list* &rarr;  Information about server disks.
    - `model` *string* &rarr;  The model of the disk.

    - `path` *string* &rarr;  The path to the disk.

    - `size` *number* &rarr;  The size of the disk.

    - `type` *string* &rarr;  The type of the disk.


- `name` *string* &rarr;  The name of the server.

- `os_type` *string* &rarr;  Server Operation System type.

- `power_state` *string* &rarr;  Server power state.

- `raid_type` *string* &rarr;  Server raid type.

- `ram_megabytes` *number* &rarr;  Server memory size in megabytes.

- `status` *string* &rarr;  Server status.

- `tags` *string* &rarr;  Server tags.

- `target_boot_order`  *list* &rarr;  Current server boot order.
    - `boot_device_type` *string* &rarr;  The boot device type.



