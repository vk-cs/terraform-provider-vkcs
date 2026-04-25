---
subcategory: "Baremetal"
layout: "vkcs"
page_title: "vkcs: vkcs_baremetal_flavor"
description: |-
  Get information on a VKCS bare metal flavor.
---

# vkcs_baremetal_flavor

Use this data source to get information about a VKCS Baremetal Flavor.

## Example Usage

```terraform
data "vkcs_baremetal_flavor" "main" {
  name              = "BM_CX301_N_BOND"
  cpu_model         = "Intel(R) Xeon(R) Gold 6338 CPU @ 2.00GHz"
  cpu_cores         = 32
  ram_size          = 128
  ssd_size          = 900
  hdd_size          = 16000
  bond_vlan_capable = true
}

output "flavor_output" {
  value = {
    id                = data.vkcs_baremetal_flavor.main.id
    name              = data.vkcs_baremetal_flavor.main.name
    cpu_cores         = data.vkcs_baremetal_flavor.main.cpu_cores
    ram_size          = data.vkcs_baremetal_flavor.main.ram_size
    ssd_size          = data.vkcs_baremetal_flavor.main.ssd_size
    hdd_size          = data.vkcs_baremetal_flavor.main.hdd_size
    bond_vlan_capable = data.vkcs_baremetal_flavor.main.bond_vlan_capable
  }
}
```

## Argument Reference
- `bond_vlan_capable` optional *boolean* &rarr;  Bond and VLAN capable.

- `cpu_cores` optional *number* &rarr;  CPU core count including hyper-threading.

- `cpu_model` optional *string* &rarr;  The CPU model.

- `hdd_size` optional *number* &rarr;  HDD size in gigabytes.

- `id` optional *string* &rarr;  The UUID of the flavor.

- `name` optional *string* &rarr;  The name of the flavor.

- `ram_size` optional *number* &rarr;  RAM in gigabytes.

- `region` optional *string* &rarr;  The region to fetch the bare metal flavor from, defaults to the provider's region.

- `ssd_size` optional *number* &rarr;  SSD size in gigabytes.


## Attributes Reference
No additional attributes are exported.

