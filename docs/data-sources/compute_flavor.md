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

```terraform
data "vkcs_compute_flavor" "small" {
  vcpus = 1
  ram   = 512
}
```

## Argument Reference
- `disk` optional *number* &rarr;  The exact amount of disk (in gigabytes).

- `flavor_id` optional *string* &rarr;  The ID of the flavor. Conflicts with the `name`, `min_ram` and `min_disk`

- `is_public` optional *boolean* &rarr;  The flavor visibility.

- `min_disk` optional *number* &rarr;  The minimum amount of disk (in gigabytes). Conflicts with the `flavor_id`.

- `min_ram` optional *number* &rarr;  The minimum amount of RAM (in megabytes). Conflicts with the `flavor_id`.

- `name` optional *string* &rarr;  The name of the flavor. Conflicts with the `flavor_id`.

- `ram` optional *number* &rarr;  The exact amount of RAM (in megabytes).

- `region` optional *string* &rarr;  The region in which to obtain the Compute client. If omitted, the `region` argument of the provider is used.

- `rx_tx_factor` optional *number* &rarr;  The `rx_tx_factor` of the flavor.

- `swap` optional *number* &rarr;  The amount of swap (in gigabytes).

- `vcpus` optional *number* &rarr;  The amount of VCPUs.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `extra_specs` *map of* *string* &rarr;  Key/Value pairs of metadata for the flavor.

- `id` *string* &rarr;  ID of the found flavor.


