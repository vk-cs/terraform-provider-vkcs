---
subcategory: "Backup"
layout: "vkcs"
page_title: "vkcs: vkcs_backup_plan"
description: |-
  Get information on an VKCS backup plan.
---

# vkcs_backup_plan

Use this data source to get backup plan info

**New since v0.4.0**.

## Example Usage

```terraform
data "vkcs_backup_plan" "plan_datasource" {
  name = vkcs_backup_plan.backup_plan.name
}
```

## Argument Reference
- `instance_id` optional *string* &rarr;  ID of the instance that should be included in backup plan

- `name` optional *string* &rarr;  Name of the backup plan

- `region` optional *string* &rarr;  The `region` to fetch availability zones from, defaults to the provider's `region`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `backup_targets`  *set* &rarr;  Set of backup targets specifying instance_id and volume_ids for each instance.<br>**New since v0.13.1**.
    - `instance_id` *string* &rarr;  ID of the instance for which specific volumes are backed up.

    - `volume_ids` *set of* *string* &rarr;  Set of volume IDs to back up for the instance. If no list is specified, backups will be created for all disks.


- `full_retention`  &rarr;  Parameters for full retention policy. Specifies number of full backups stored. Incremental backups (if enabled) are not counted as full
    - `max_full_backup` *number* &rarr;  Maximum number of backups


- `gfs_retention`  &rarr;  Parameters for gfs retention policy. Specifies number of full backups stored. Incremental backups (if enabled) are not counted as full
    - `gfs_monthly` *number* &rarr;  Number of months to store backups

    - `gfs_weekly` *number* &rarr;  Number of weeks to store backups

    - `gfs_yearly` *number* &rarr;  Number of years to store backups


- `id` *string* &rarr;  ID of the resource

- `incremental_backup` *boolean* &rarr;  Whether incremental backups should be stored

- `instance_ids` *set of* *string* &rarr;  Set of ids of backed up instances

- `provider_id` *string* &rarr;  ID of backup provider

- `schedule` 
    - `date` *string* &rarr;  List of days when to perform backups. If incremental_backups is enabled, this field contains day of full backup

    - `every_hours` *number* &rarr;  Hour period of backups

    - `time` *string* &rarr;  Time of backup in format hh:mm, using UTC timezone



