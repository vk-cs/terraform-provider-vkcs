---
subcategory: "Baremetal"
layout: "vkcs"
page_title: "vkcs: vkcs_baremetal_flavors"
description: |-
  Get information on VKCS bare metal flavors.
---

# vkcs_baremetal_flavors

Use this data source to get a list of available VKCS Baremetal Flavors.

## Example Usage

```terraform
data "vkcs_baremetal_flavors" "main" {}

output "flavors_output" {
  value = {
    flavors = data.vkcs_baremetal_flavors.main
  }
}
```

## Argument Reference
- `region` optional *string* &rarr;  The region to fetch the bare metal flavor from, defaults to the provider's region.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `flavors`  *list* &rarr;  Available Baremetal Flavors.
    - `bond_vlan_capable` *boolean* &rarr;  Bond and VLAN capable.

    - `cpu_cores` *number* &rarr;  CPU core count including hyper-threading.

    - `cpu_model` *string* &rarr;  The CPU model.

    - `hdd_size` *number* &rarr;  HDD size in gigabytes.

    - `name` *string* &rarr;  The name of the flavor.

    - `ram_size` *number* &rarr;  RAM in gigabytes.

    - `ssd_size` *number* &rarr;  SSD size in gigabytes.



