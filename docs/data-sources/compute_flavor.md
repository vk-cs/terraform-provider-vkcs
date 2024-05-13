---
subcategory: "Virtual Machines"
layout: "vkcs"
page_title: "vkcs: vkcs_compute_flavor"
description: |-
  Get information on an VKCS Flavor.
---

# vkcs_compute_flavor

Use this data source to get the ID of an available VKCS flavor.

## Example Usage
### Filter by name
```terraform
data "vkcs_compute_flavor" "basic" {
  name = "Standard-2-8-50"
}
```

### Filter by number of vCPUs, RAM and CPU generation
```terraform
data "vkcs_compute_flavor" "basic" {
  vcpus = 1
  ram   = 1024
  # specify cpu_generation to distinguish between several flavors with the same CPU and RAM 
  extra_specs = {
    "mcs:cpu_generation" : "cascadelake-v1"
  }
}
```

### Filter by number of vCPUs and minimum RAM
```terraform
# If the exact amount of RAM is not so important to you, then you can specify the minimum value that will satisfy you 
# and flavor with minimum of ram will be automatically selected for you.
data "vkcs_compute_flavor" "standard_4_min_6gb" {
  vcpus   = 4
  min_ram = 6000
}
```

## Argument Reference
- `disk` optional *number* &rarr;  The exact amount of disk (in gigabytes). Don't set disk, when min_disk is set.

- `extra_specs` optional *map of* *string* &rarr;  Key/Value pairs of metadata for the flavor. Be careful when using it, there is no validation applied to this field. When searching for a suitable flavor, it checks all required extra specs in a flavor metadata. See https://cloud.vk.com/docs/base/iaas/concepts/vm-concept

- `flavor_id` optional *string* &rarr;  The ID of the flavor. Conflicts with the `name`, `min_ram` and `min_disk`

- `is_public` optional *boolean* &rarr;  The flavor visibility.

- `min_disk` optional *number* &rarr;  The minimum amount of disk (in gigabytes). Conflicts with the `flavor_id`.

- `min_ram` optional *number* &rarr;  The minimum amount of RAM (in megabytes). Conflicts with the `flavor_id`.

- `name` optional *string* &rarr;  The name of the flavor. Conflicts with the `flavor_id`.

- `ram` optional *number* &rarr;  The exact amount of RAM (in megabytes). Don't set ram, when min_ram is set.

- `region` optional *string* &rarr;  The region in which to obtain the Compute client. If omitted, the `region` argument of the provider is used.

- `rx_tx_factor` optional *number* &rarr;  The `rx_tx_factor` of the flavor.

- `swap` optional *number* &rarr;  The amount of swap (in gigabytes).

- `vcpus` optional *number* &rarr;  The amount of VCPUs.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the found flavor.


