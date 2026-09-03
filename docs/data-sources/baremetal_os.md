---
subcategory: "Baremetal"
layout: "vkcs"
page_title: "vkcs: vkcs_baremetal_os"
description: |-
  Get information on a VKCS bare metal operating system.
---

# vkcs_baremetal_os

Use this data source to get information about a VKCS baremetal OS.

## Example Usage

```terraform
data "vkcs_baremetal_os" "ubuntu" {
  name      = "ubuntu"
  version   = "24.04"
  raid_type = "RAID1"
}

output "flavor_output" {
  value = {
    id        = data.vkcs_baremetal_os.ubuntu.id
    name      = data.vkcs_baremetal_os.ubuntu.name
    version   = data.vkcs_baremetal_os.ubuntu.version
    raid_type = data.vkcs_baremetal_os.ubuntu.raid_type
  }
}
```

## Argument Reference
- `id` optional *string* &rarr;  The UUID of the OS.

- `name` optional *string* &rarr;  The name of the OS.

- `raid_type` optional *string* &rarr;  The raid type of the OS.

- `region` optional *string* &rarr;  The region to fetch the bare metal OS from, defaults to the provider's region.

- `version` optional *string* &rarr;  The version of the OS.


## Attributes Reference
No additional attributes are exported.

