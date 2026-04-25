---
subcategory: "Baremetal"
layout: "vkcs"
page_title: "vkcs: vkcs_baremetal_oses"
description: |-
  Get information on VKCS bare metal operating systems.
---

# vkcs_baremetal_oses

Use this data source to get a list of available VKCS Baremetal OSes.

## Example Usage

```terraform
data "vkcs_baremetal_oses" "main" {}

output "oses_output" {
  value = {
    oses = data.vkcs_baremetal_oses.main
  }
}
```

## Argument Reference
- `region` optional *string* &rarr;  The region to fetch the bare metal OSes from, defaults to the provider's region.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `oses`  *list* &rarr;  Available Baremetal OSes.
    - `name` *string* &rarr;  The name of the OS.

    - `raid_type` *string* &rarr;  The raid type of the OS.

    - `version` *string* &rarr;  The version of the OS.



