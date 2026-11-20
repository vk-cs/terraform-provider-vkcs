---
subcategory: "Backup"
layout: "vkcs"
page_title: "vkcs: vkcs_backup_checkpoints"
description: |-
  Reads backup checkpoints and their Object Lock settings for a VKCS backup plan.
---

# vkcs_backup_checkpoints

Use this data source to read backup checkpoints grouped by plan. The object_lock values are read-only.

## Example Usage

```terraform
data "vkcs_backup_checkpoints" "plan" {
  count   = var.enable_checkpoints ? 1 : 0
  plan_id = var.backup_plan_id
}

output "legal_hold_checkpoint_ids" {
  value = flatten([
    for resource in try(data.vkcs_backup_checkpoints.plan[0].resources, []) : [
      for backup in resource.backups : backup.checkpoint_id
      if try(backup.object_lock.type, "") == "legalhold"
    ]
  ])
}
```
```terraform
variable "enable_checkpoints" {
  type        = bool
  default     = false
  description = "Set to true to query checkpoints from an existing backup plan."
}

variable "backup_plan_id" {
  type        = string
  default     = ""
  description = "ID of an existing backup plan; used when enable_checkpoints is true."
}
```

## Argument Reference
- `plan_id` **required** *string* &rarr;  ID of the backup plan whose checkpoints should be retrieved.

- `region` optional *string* &rarr;  The region in which the backup plan exists. Defaults to the provider region.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the backup plan used to query checkpoints.

- `project_id` *string* &rarr;  ID of the project containing the plan and its checkpoints.

- `resources`  *list* &rarr;  Resources protected by the plan, with their backup checkpoints.
    - `backups`  *list* &rarr;  Backup checkpoints created for this resource.
        - `checkpoint_id` *string* &rarr;  ID of the backup checkpoint.

        - `created_at` *string* &rarr;  Checkpoint creation timestamp.

        - `object_lock`  &rarr;  Object Lock protection for the checkpoint; read-only.
            - `days` *number* &rarr;  Retention period in days, if configured.

            - `retain_until` *number* &rarr;  Retention end date as a Unix timestamp, if configured.

            - `type` *string* &rarr;  Object Lock mode: compliance, governance, or legalhold.


        - `status` *string* &rarr;  Status of the backup checkpoint.


    - `id` *string* &rarr;  ID of the protected resource.

    - `name` *string* &rarr;  Name of the protected resource.

    - `resource_type` *string* &rarr;  Type of the protected resource.
