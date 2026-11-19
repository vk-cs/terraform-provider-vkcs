---
subcategory: "Backup"
layout: "vkcs"
page_title: "vkcs: vkcs_backup_restore"
description: |-
  Reads restore details and per-resource statuses for a VKCS backup restore.
---

# vkcs_backup_restore

Use this data source to read a restore operation and its per-resource statuses.

Use the computed `restore_id` from `vkcs_compute_instance` to inspect the operation created from `checkpoint_id`.

## Example Usage

```terraform
resource "vkcs_compute_instance" "restored_vm" {
  count              = var.enable_restore ? 1 : 0
  name               = var.instance_name
  checkpoint_id      = var.checkpoint_id
  source_resource_id = var.source_resource_id != "" ? var.source_resource_id : null
  key_pair           = var.key_pair

  network {
    uuid = var.network_id
  }
}

data "vkcs_backup_restore" "status" {
  count      = var.enable_restore ? 1 : 0
  restore_id = try(vkcs_compute_instance.restored_vm[0].restore_id, "")
}
```
```terraform
variable "enable_restore" {
  type        = bool
  default     = false
  description = "Set to true to restore a VM from the selected checkpoint."
}

variable "instance_name" {
  type        = string
  default     = "restored-vm-example"
  description = "Name for the restored VM."
}

variable "checkpoint_id" {
  type        = string
  default     = ""
  description = "ID of the backup checkpoint to restore."
}

variable "source_resource_id" {
  type        = string
  default     = ""
  description = "Nova server resource ID within the checkpoint when it contains multiple VMs."
}

variable "network_id" {
  type        = string
  default     = ""
  description = "Network ID for the restored VM."
}

variable "key_pair" {
  type        = string
  default     = ""
  description = "Keypair name for the restored VM."
}
```

## Argument Reference
- `restore_id` **required** *string* &rarr;  ID of the restore operation to retrieve.

- `region` optional *string* &rarr;  The region from which to obtain the backup service endpoint. Defaults to the provider region.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `checkpoint_id` *string* &rarr;  ID of the restored backup checkpoint.

- `id` *string* &rarr;  ID of the restore operation.

- `provider_id` *string* &rarr;  ID of the backup provider used for restoration.

- `resources_reason` *map of* *string* &rarr;  Reason for each resource that could not be restored.

- `resources_status` *map of* *string* &rarr;  Restore status for each affected resource.
