---
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
- `disk` **Number** (*Optional*) The exact amount of disk (in gigabytes).

- `flavor_id` **String** (*Optional*) The ID of the flavor. Conflicts with the `name`, `min_ram` and `min_disk`

- `is_public` **Boolean** (*Optional*) The flavor visibility.

- `min_disk` **Number** (*Optional*) The minimum amount of disk (in gigabytes). Conflicts with the `flavor_id`.

- `min_ram` **Number** (*Optional*) The minimum amount of RAM (in megabytes). Conflicts with the `flavor_id`.

- `name` **String** (*Optional*) The name of the flavor. Conflicts with the `flavor_id`.

- `ram` **Number** (*Optional*) The exact amount of RAM (in megabytes).

- `region` **String** (*Optional*) The region in which to obtain the Compute client. If omitted, the `region` argument of the provider is used.

- `rx_tx_factor` **Number** (*Optional*) The `rx_tx_factor` of the flavor.

- `swap` **Number** (*Optional*) The amount of swap (in gigabytes).

- `vcpus` **Number** (*Optional*) The amount of VCPUs.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `extra_specs` <strong>Map of </strong>**String** Key/Value pairs of metadata for the flavor.

- `id` **String** ID of the found flavor.


