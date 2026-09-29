---
subcategory: "Baremetal"
layout: "vkcs"
page_title: "vkcs: vkcs_baremetal_flavor_display_name"
description: |-
  Manages the project-specific display name of an existing VKCS bare metal flavor.
---

# vkcs_baremetal_flavor_display_name

Manages the project-specific display name of an existing VKCS bare metal flavor. The flavor itself is not created or deleted.

## Example Usage

```terraform
data "vkcs_baremetal_flavor" "selected" {
  name = "BM_CX301_N_BOND"
}

# Get the current project-specific display name.
output "display_name_before" {
  value = data.vkcs_baremetal_flavor.selected.display_name
}

# Write a project-specific display name.
resource "vkcs_baremetal_flavor_display_name" "main" {
  id           = data.vkcs_baremetal_flavor.selected.id
  display_name = "Terraform managed bare metal flavor"
}

# Rewrite the display name by changing the value above and running:
# terraform apply

output "display_name" {
  value = vkcs_baremetal_flavor_display_name.main.display_name
}

# Reset the display name by running:
# terraform destroy
# Destroying this resource clears the project-specific display name.
```

## Argument Reference
- `display_name` **required** *string* &rarr;  The project-specific display name of the flavor.

- `id` **required** *string* &rarr;  The UUID of the existing flavor.


## Attributes Reference
No additional attributes are exported.


Destroying the resource clears the project-specific display name. The flavor is not deleted.

## Import

An existing flavor can be imported using its ID:

```shell
terraform import vkcs_baremetal_flavor_display_name.main <flavor-id>
```
